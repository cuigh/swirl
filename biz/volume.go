package biz

import (
	"context"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
)

type VolumeBiz interface {
	Search(ctx context.Context, node, name string, pageIndex, pageSize int) ([]*Volume, int, error)
	Find(ctx context.Context, node, name string) (volume *Volume, raw string, err error)
	Delete(ctx context.Context, node, name string, user web.User) (err error)
	Create(ctx context.Context, volume *Volume, user web.User) (err error)
	Prune(ctx context.Context, node string, user web.User) (count int, size uint64, err error)
}

func NewVolume(d *docker.Docker, eb EventBiz) VolumeBiz {
	return &volumeBiz{d: d, eb: eb}
}

type volumeBiz struct {
	d  *docker.Docker
	eb EventBiz
}

func (b *volumeBiz) Find(ctx context.Context, node, name string) (volume *Volume, raw string, err error) {
	var (
		v types.Volume
		r []byte
	)

	if v, r, err = b.d.VolumeInspect(ctx, node, name); err == nil {
		raw, err = indentJSON(r)
	}

	if err == nil {
		volume = newVolume(&v)
	}
	return
}

func (b *volumeBiz) Search(ctx context.Context, node, name string, pageIndex, pageSize int) (volumes []*Volume, total int, err error) {
	list, total, err := b.d.VolumeList(ctx, node, name, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	volumes = make([]*Volume, len(list))
	for i, v := range list {
		volumes[i] = newVolume(v)
	}
	return volumes, total, nil
}

func (b *volumeBiz) Delete(ctx context.Context, node, name string, user web.User) (err error) {
	err = b.d.VolumeRemove(ctx, node, name)
	if err == nil {
		b.eb.CreateVolume(EventActionDelete, node, name, user)
	}
	return
}

func (b *volumeBiz) Create(ctx context.Context, vol *Volume, user web.User) (err error) {
	options := &volume.VolumeCreateBody{
		Name:       vol.Name,
		Driver:     vol.Driver,
		DriverOpts: toMap(vol.Options),
		Labels:     toMap(vol.Labels),
	}
	if vol.Driver == "other" {
		options.Driver = vol.CustomDriver
	} else {
		options.Driver = vol.Driver
	}

	err = b.d.VolumeCreate(ctx, vol.Node, options)
	if err == nil {
		b.eb.CreateVolume(EventActionDelete, vol.Node, vol.Name, user)
	}
	return
}

func (b *volumeBiz) Prune(ctx context.Context, node string, user web.User) (count int, size uint64, err error) {
	var report types.VolumesPruneReport
	report, err = b.d.VolumePrune(ctx, node)
	if err == nil {
		count, size = len(report.VolumesDeleted), report.SpaceReclaimed
		b.eb.CreateVolume(EventActionPrune, node, "", user)
	}
	return
}

type Volume struct {
	Node         string                 `json:"node"`
	Name         string                 `json:"name"`
	Driver       string                 `json:"driver,omitempty"`
	CustomDriver string                 `json:"customDriver,omitempty"`
	CreatedAt    string                 `json:"createdAt"`
	MountPoint   string                 `json:"mountPoint,omitempty"`
	Scope        string                 `json:"scope"`
	Labels       data.Options           `json:"labels,omitempty"`
	Options      data.Options           `json:"options,omitempty"`
	Status       map[string]interface{} `json:"status,omitempty"`
	RefCount     int64                  `json:"refCount"`
	Size         int64                  `json:"size"`
}

func newVolume(v *types.Volume) *Volume {
	createdAt, _ := time.Parse(time.RFC3339Nano, v.CreatedAt)
	vol := &Volume{
		Name:       v.Name,
		Driver:     v.Driver,
		CreatedAt:  formatTime(createdAt),
		MountPoint: v.Mountpoint,
		Scope:      v.Scope,
		Status:     v.Status,
		Labels:     mapToOptions(v.Labels),
		Options:    mapToOptions(v.Options),
		RefCount:   -1,
		Size:       -1,
	}
	if v.UsageData != nil {
		vol.RefCount = v.UsageData.RefCount
		vol.Size = v.UsageData.Size
	}
	return vol
}

package biz

import (
	"context"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
)

type ImageBiz interface {
	Search(node, name string, pageIndex, pageSize int) ([]*Image, int, error)
	Find(node, name string) (image *Image, raw string, err error)
	Delete(node, id string, user web.User) (err error)
	Prune(node string, user web.User) (count int, size uint64, err error)
}

func NewImage(d *docker.Docker, eb EventBiz) ImageBiz {
	return &imageBiz{d: d, eb: eb}
}

type imageBiz struct {
	d  *docker.Docker
	eb EventBiz
}

func (b *imageBiz) Find(node, id string) (img *Image, raw string, err error) {
	var (
		i         types.ImageInspect
		r         []byte
		histories []image.HistoryResponseItem
		ctx       = context.TODO()
	)

	if i, r, err = b.d.ImageInspect(ctx, node, id); err == nil {
		raw, err = indentJSON(r)
	}

	if err == nil {
		histories, err = b.d.ImageHistory(ctx, node, id)
	}

	if err == nil {
		img = newImageDetail(&i, histories)
	}
	return
}

func (b *imageBiz) Search(node, name string, pageIndex, pageSize int) (images []*Image, total int, err error) {
	list, total, err := b.d.ImageList(context.TODO(), node, name, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	images = make([]*Image, len(list))
	for i, nr := range list {
		images[i] = newImageSummary(&nr)
	}
	return images, total, nil
}

func (b *imageBiz) Delete(node, id string, user web.User) (err error) {
	err = b.d.ImageRemove(context.TODO(), node, id)
	if err == nil {
		b.eb.CreateImage(EventActionDelete, id, user)
	}
	return
}

func (b *imageBiz) Prune(node string, user web.User) (count int, size uint64, err error) {
	var report types.ImagesPruneReport
	if report, err = b.d.ImagePrune(context.TODO(), node); err == nil {
		count, size = len(report.ImagesDeleted), report.SpaceReclaimed
		b.eb.CreateImage(EventActionPrune, "", user)
	}
	return
}

type Image struct {
	/* Summary */
	ID          string       `json:"id"`
	ParentID    string       `json:"pid,omitempty"`
	Created     string       `json:"created"`
	Containers  int64        `json:"containers"`
	Digests     []string     `json:"digests"`
	Tags        []string     `json:"tags"`
	Labels      data.Options `json:"labels"`
	Size        int64        `json:"size"`
	SharedSize  int64        `json:"sharedSize"`
	VirtualSize int64        `json:"virtualSize"`

	/* Detail */
	Comment       string           `json:"comment,omitempty"`
	Container     string           `json:"container,omitempty"`
	DockerVersion string           `json:"dockerVersion,omitempty"`
	Author        string           `json:"author,omitempty"`
	Architecture  string           `json:"arch,omitempty"`
	Variant       string           `json:"variant,omitempty"`
	OS            string           `json:"os,omitempty"`
	OSVersion     string           `json:"osVersion,omitempty"`
	GraphDriver   ImageGraphDriver `json:"graphDriver"`
	RootFS        ImageRootFS      `json:"rootFS"`
	LastTagTime   string           `json:"lastTagTime,omitempty"`
	Histories     []*ImageHistory  `json:"histories,omitempty"`
	//Config          *container.Config
	//ContainerConfig *container.Config
}

type ImageGraphDriver struct {
	Name string       `json:"name,omitempty"`
	Data data.Options `json:"data,omitempty"`
}

type ImageRootFS struct {
	Type      string   `json:"type"`
	Layers    []string `json:"layers,omitempty"`
	BaseLayer string   `json:"baseLayer,omitempty"`
}

type ImageHistory struct {
	ID        string   `json:"id,omitempty"`
	Comment   string   `json:"comment,omitempty"`
	Size      int64    `json:"size,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
	CreatedBy string   `json:"createdBy,omitempty"`
}

func newImageSummary(is *types.ImageSummary) *Image {
	i := &Image{
		ID:          is.ID,
		ParentID:    is.ParentID,
		Created:     formatTime(time.Unix(is.Created, 0)),
		Containers:  is.Containers,
		Digests:     is.RepoDigests,
		Tags:        is.RepoTags,
		Labels:      mapToOptions(is.Labels),
		SharedSize:  is.SharedSize,
		Size:        is.Size,
		VirtualSize: is.VirtualSize,
	}
	return i
}

func newImageDetail(is *types.ImageInspect, items []image.HistoryResponseItem) *Image {
	created, _ := time.Parse(time.RFC3339Nano, is.Created)
	histories := make([]*ImageHistory, len(items))
	for i, item := range items {
		histories[i] = &ImageHistory{
			ID:        item.ID,
			Comment:   item.Comment,
			Size:      item.Size,
			Tags:      item.Tags,
			CreatedAt: formatTime(time.Unix(item.Created, 0)),
			CreatedBy: item.CreatedBy,
		}
	}

	i := &Image{
		ID:       is.ID,
		ParentID: is.Parent,
		Created:  formatTime(created),
		Digests:  is.RepoDigests,
		Tags:     is.RepoTags,
		//Labels:      mapToOptions(is.Labels),
		Size:        is.Size,
		VirtualSize: is.VirtualSize,

		Comment:       is.Comment,
		Container:     is.Container,
		DockerVersion: is.DockerVersion,
		Author:        is.Author,
		Architecture:  is.Architecture,
		Variant:       is.Variant,
		OS:            is.Os,
		OSVersion:     is.OsVersion,
		LastTagTime:   formatTime(is.Metadata.LastTagTime),
		GraphDriver: ImageGraphDriver{
			Name: is.GraphDriver.Name,
			Data: mapToOptions(is.GraphDriver.Data),
		},
		RootFS: ImageRootFS{
			Type:      is.RootFS.Type,
			Layers:    is.RootFS.Layers,
			BaseLayer: is.RootFS.BaseLayer,
		},
		Histories: histories,
	}
	return i
}

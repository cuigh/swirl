package docker

import (
	"context"
	"sort"

	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// VolumeList return volumes on the host.
func VolumeList(name string, pageIndex, pageSize int) (volumes []*types.Volume, totalCount int, err error) {
	var (
		ctx  context.Context
		cli  *client.Client
		resp volume.VolumesListOKBody
	)

	ctx, cli, err = mgr.Client()
	if err != nil {
		return
	}

	f := filters.NewArgs()
	//f.Add("dangling", "true")
	//f.Add("driver", "xx")
	if name != "" {
		f.Add("name", name)
	}
	resp, err = cli.VolumeList(ctx, f)
	if err != nil {
		return
	}
	sort.Slice(resp.Volumes, func(i, j int) bool {
		return resp.Volumes[i].Name < resp.Volumes[j].Name
	})

	totalCount = len(resp.Volumes)
	start, end := misc.Page(totalCount, pageIndex, pageSize)
	volumes = resp.Volumes[start:end]
	return
}

// VolumeCreate create a volume.
func VolumeCreate(info *model.VolumeCreateInfo) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		options := volume.VolumesCreateBody{
			Name:       info.Name,
			DriverOpts: info.Options.ToMap(),
			Labels:     info.Labels.ToMap(),
		}
		if info.Driver == "other" {
			options.Driver = info.CustomDriver
		} else {
			options.Driver = info.Driver
		}

		_, err = cli.VolumeCreate(ctx, options)
		return
	})
}

// VolumeRemove remove a volume.
func VolumeRemove(name string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		return cli.VolumeRemove(ctx, name, false)
	})
}

// VolumePrune remove all unused volumes.
func VolumePrune() (report types.VolumesPruneReport, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		f := filters.NewArgs()
		report, err = cli.VolumesPrune(ctx, f)
		return
	})
	return
}

// VolumeInspectRaw return volume raw information.
func VolumeInspectRaw(name string) (vol types.Volume, raw []byte, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) error {
		vol, raw, err = cli.VolumeInspectWithRaw(ctx, name)
		return err
	})
	return
}

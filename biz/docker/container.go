package docker

import (
	"context"

	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// ContainerList return containers on the host.
func ContainerList(name string, pageIndex, pageSize int) (infos []*model.ContainerListInfo, totalCount int, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		var (
			containers []types.Container
			opts       = types.ContainerListOptions{}
		)

		if name != "" {
			opts.Filters = filters.NewArgs()
			opts.Filters.Add("name", name)
		}
		containers, err = cli.ContainerList(ctx, opts)
		if err == nil {
			//sort.Slice(containers, func(i, j int) bool {
			//	return containers[i] < containers[j].Description.Hostname
			//})
			totalCount = len(containers)
			start, end := misc.Page(totalCount, pageIndex, pageSize)
			containers = containers[start:end]
			if length := len(containers); length > 0 {
				infos = make([]*model.ContainerListInfo, length)
				for i, c := range containers {
					infos[i] = model.NewContainerListInfo(c)
				}
			}
		}
		return
	})
	return
}

// ContainerInspect return detail information of a container.
func ContainerInspect(id string) (container types.ContainerJSON, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		container, err = cli.ContainerInspect(ctx, id)
		return
	})
	return
}

// ContainerInspectRaw return container raw information.
func ContainerInspectRaw(id string) (container types.ContainerJSON, raw []byte, err error) {
	mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		container, raw, err = cli.ContainerInspectWithRaw(ctx, id, true)
		return
	})
	return
}

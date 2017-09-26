package docker

import (
	"context"
	"sort"

	"github.com/cuigh/swirl/misc"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// ConfigList return all configs.
func ConfigList(name string, pageIndex, pageSize int) (configs []swarm.Config, totalCount int, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		opts := types.ConfigListOptions{}
		if name != "" {
			opts.Filters = filters.NewArgs()
			opts.Filters.Add("name", name)
		}
		configs, err = cli.ConfigList(ctx, opts)
		if err == nil {
			sort.Slice(configs, func(i, j int) bool {
				return configs[i].Spec.Name < configs[j].Spec.Name
			})
			totalCount = len(configs)
			start, end := misc.Page(totalCount, pageIndex, pageSize)
			configs = configs[start:end]
		}
		return
	})
	return
}

// ConfigCreate create a config.
func ConfigCreate(name string, data []byte, labels map[string]string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		// todo:
		spec := swarm.ConfigSpec{}
		spec.Name = name
		spec.Data = data
		spec.Labels = labels
		_, err = cli.ConfigCreate(ctx, spec)
		return
	})
}

// ConfigRemove remove a config.
func ConfigRemove(ids []string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		for _, id := range ids {
			if err = cli.ConfigRemove(ctx, id); err != nil {
				return
			}
		}
		return
	})
}

package docker

import (
	"context"
	"sort"

	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
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
func ConfigCreate(info *model.ConfigCreateInfo) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		spec := swarm.ConfigSpec{}
		spec.Name = info.Name
		spec.Data = []byte(info.Data)
		spec.Labels = info.Labels.ToMap()
		_, err = cli.ConfigCreate(ctx, spec)
		return
	})
}

// ConfigUpdate update a config.
func ConfigUpdate(info *model.ConfigUpdateInfo) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		var cfg swarm.Config
		cfg, _, err = cli.ConfigInspectWithRaw(ctx, info.ID)
		if err != nil {
			return err
		}

		spec := cfg.Spec
		// only the Labels field can be updated on API 1.30
		//spec.Name = info.Name
		//spec.Data = []byte(info.Data)
		spec.Labels = info.Labels.ToMap()
		return cli.ConfigUpdate(ctx, info.ID, version(info.Version), spec)
	})
}

// ConfigInspect returns config information with raw data.
func ConfigInspect(id string) (cfg swarm.Config, raw []byte, err error) {
	var (
		ctx context.Context
		cli *client.Client
	)
	if ctx, cli, err = mgr.Client(); err == nil {
		cfg, raw, err = cli.ConfigInspectWithRaw(ctx, id)
	}
	return
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

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
func (d *Docker) ConfigList(ctx context.Context, name string, pageIndex, pageSize int) (configs []swarm.Config, total int, err error) {
	err = d.call(func(c *client.Client) (err error) {
		opts := types.ConfigListOptions{}
		if name != "" {
			opts.Filters = filters.NewArgs()
			opts.Filters.Add("name", name)
		}
		configs, err = c.ConfigList(ctx, opts)
		if err == nil {
			sort.Slice(configs, func(i, j int) bool {
				return configs[i].Spec.Name < configs[j].Spec.Name
			})
			total = len(configs)
			start, end := misc.Page(total, pageIndex, pageSize)
			configs = configs[start:end]
		}
		return
	})
	return
}

// ConfigInspect returns config information with raw data.
func (d *Docker) ConfigInspect(ctx context.Context, id string) (cfg swarm.Config, raw []byte, err error) {
	err = d.call(func(c *client.Client) (err error) {
		cfg, raw, err = c.ConfigInspectWithRaw(ctx, id)
		return
	})
	return
}

// ConfigRemove remove a config.
func (d *Docker) ConfigRemove(ctx context.Context, id string) error {
	return d.call(func(c *client.Client) (err error) {
		return c.ConfigRemove(ctx, id)
	})
}

// ConfigCreate create a config.
func (d *Docker) ConfigCreate(ctx context.Context, spec *swarm.ConfigSpec) (id string, err error) {
	err = d.call(func(c *client.Client) (err error) {
		var resp types.ConfigCreateResponse
		if resp, err = c.ConfigCreate(ctx, *spec); err == nil {
			id = resp.ID
		}
		return
	})
	return
}

// ConfigUpdate update a config.
func (d *Docker) ConfigUpdate(ctx context.Context, id string, version uint64, spec *swarm.ConfigSpec) error {
	return d.call(func(c *client.Client) (err error) {
		return c.ConfigUpdate(ctx, id, newVersion(version), *spec)
	})
}

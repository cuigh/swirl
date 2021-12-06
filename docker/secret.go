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

// SecretList return all secrets.
func (d *Docker) SecretList(ctx context.Context, name string, pageIndex, pageSize int) (secrets []swarm.Secret, total int, err error) {
	err = d.call(func(c *client.Client) (err error) {
		opts := types.SecretListOptions{}
		if name != "" {
			opts.Filters = filters.NewArgs()
			opts.Filters.Add("name", name)
		}
		secrets, err = c.SecretList(ctx, opts)
		if err == nil {
			sort.Slice(secrets, func(i, j int) bool {
				return secrets[i].Spec.Name < secrets[j].Spec.Name
			})
			total = len(secrets)
			start, end := misc.Page(total, pageIndex, pageSize)
			secrets = secrets[start:end]
		}
		return
	})
	return
}

// SecretInspect returns secret information with raw data.
func (d *Docker) SecretInspect(ctx context.Context, id string) (secret swarm.Secret, raw []byte, err error) {
	err = d.call(func(c *client.Client) (err error) {
		secret, raw, err = c.SecretInspectWithRaw(ctx, id)
		return
	})
	return
}

// SecretRemove remove a secret.
func (d *Docker) SecretRemove(ctx context.Context, id string) error {
	return d.call(func(c *client.Client) (err error) {
		return c.SecretRemove(ctx, id)
	})
}

// SecretCreate create a secret.
func (d *Docker) SecretCreate(ctx context.Context, spec *swarm.SecretSpec) (id string, err error) {
	err = d.call(func(c *client.Client) (err error) {
		var resp types.SecretCreateResponse
		if resp, err = c.SecretCreate(ctx, *spec); err == nil {
			id = resp.ID
		}
		return
	})
	return
}

// SecretUpdate update a config.
func (d *Docker) SecretUpdate(ctx context.Context, id string, version uint64, spec *swarm.SecretSpec) error {
	return d.call(func(c *client.Client) (err error) {
		return c.SecretUpdate(ctx, id, newVersion(version), *spec)
	})
}

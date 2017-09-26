package docker

import (
	"context"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/cuigh/swirl/misc"
)

// SecretList return all secrets.
func SecretList(name string, pageIndex, pageSize int) (secrets []swarm.Secret, totalCount int, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		opts := types.SecretListOptions{}
		if name != "" {
			opts.Filters = filters.NewArgs()
			opts.Filters.Add("name", name)
		}
		secrets, err = cli.SecretList(ctx, opts)
		if err == nil {
			sort.Slice(secrets, func(i, j int) bool {
				return secrets[i].Spec.Name < secrets[j].Spec.Name
			})
			totalCount = len(secrets)
			start, end := misc.Page(totalCount, pageIndex, pageSize)
			secrets = secrets[start:end]
		}
		return
	})
	return
}

// SecretCreate create a secret.
func SecretCreate(name string, data []byte, labels map[string]string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		// todo:
		spec := swarm.SecretSpec{}
		spec.Name = name
		spec.Data = data
		spec.Labels = labels
		_, err = cli.SecretCreate(ctx, spec)
		return
	})
}

// SecretRemove remove a secret.
func SecretRemove(id string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		return cli.SecretRemove(ctx, id)
	})
}

package docker

import (
	"context"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// NetworkCreate create a network.
func (d *Docker) NetworkCreate(ctx context.Context, name string, options *types.NetworkCreate) error {
	return d.call(func(client *client.Client) error {
		resp, err := client.NetworkCreate(ctx, name, *options)
		if err == nil && resp.Warning != "" {
			d.logger.Warnf("network '%s' was created but got warning: %s", name, resp.Warning)
		}
		return err
	})
}

// NetworkList return all networks.
func (d *Docker) NetworkList(ctx context.Context) (networks []types.NetworkResource, err error) {
	err = d.call(func(c *client.Client) (err error) {
		networks, err = c.NetworkList(ctx, types.NetworkListOptions{})
		if err == nil {
			sort.Slice(networks, func(i, j int) bool {
				return networks[i].Name < networks[j].Name
			})
		}
		return
	})
	return
}

// NetworkCount return number of networks.
func (d *Docker) NetworkCount(ctx context.Context) (count int, err error) {
	err = d.call(func(c *client.Client) (err error) {
		var networks []types.NetworkResource
		networks, err = c.NetworkList(ctx, types.NetworkListOptions{})
		if err == nil {
			count = len(networks)
		}
		return
	})
	return
}

// NetworkRemove remove a network.
func (d *Docker) NetworkRemove(ctx context.Context, name string) error {
	return d.call(func(c *client.Client) (err error) {
		return c.NetworkRemove(ctx, name)
	})
}

// NetworkDisconnect Disconnect a container from a network.
func (d *Docker) NetworkDisconnect(ctx context.Context, network, container string) error {
	return d.call(func(c *client.Client) (err error) {
		return c.NetworkDisconnect(ctx, network, container, false)
	})
}

// NetworkInspect return network information.
func (d *Docker) NetworkInspect(ctx context.Context, name string) (network types.NetworkResource, raw []byte, err error) {
	var c *client.Client
	if c, err = d.client(); err == nil {
		network, raw, err = c.NetworkInspectWithRaw(ctx, name, types.NetworkInspectOptions{})
	}
	return
}

// NetworkNames return network names by id list.
func (d *Docker) NetworkNames(ctx context.Context, ids ...string) (names map[string]string, err error) {
	var c *client.Client
	if c, err = d.client(); err == nil {
		names = make(map[string]string)
		for _, id := range ids {
			var n types.NetworkResource
			n, err = c.NetworkInspect(ctx, id, types.NetworkInspectOptions{})
			if err != nil {
				break
			}
			names[id] = n.Name
		}
	}
	return
}

package docker

import (
	"context"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// NodeList return all swarm nodes.
func (d *Docker) NodeList(ctx context.Context) (nodes []swarm.Node, err error) {
	var c *client.Client
	c, err = d.client()
	if err == nil {
		nodes, err = c.NodeList(ctx, types.NodeListOptions{})
		if err == nil {
			sort.Slice(nodes, func(i, j int) bool {
				return nodes[i].Description.Hostname < nodes[j].Description.Hostname
			})
		}
	}
	return
}

// NodeUpdate update a node.
func (d *Docker) NodeUpdate(ctx context.Context, id string, version uint64, spec *swarm.NodeSpec) error {
	return d.call(func(cli *client.Client) (err error) {
		return cli.NodeUpdate(ctx, id, newVersion(version), *spec)
	})
}

// NodeRemove remove a swarm node from cluster.
func (d *Docker) NodeRemove(ctx context.Context, id string) error {
	return d.call(func(cli *client.Client) (err error) {
		return cli.NodeRemove(ctx, id, types.NodeRemoveOptions{})
	})
}

// NodeCount return number of swarm nodes.
func (d *Docker) NodeCount(ctx context.Context) (count int, err error) {
	err = d.call(func(cli *client.Client) (err error) {
		var nodes []swarm.Node
		nodes, err = cli.NodeList(ctx, types.NodeListOptions{})
		if err == nil {
			count = len(nodes)
		}
		return
	})
	return
}

// NodeInspect return node information.
func (d *Docker) NodeInspect(ctx context.Context, id string) (node swarm.Node, raw []byte, err error) {
	err = d.call(func(cli *client.Client) (err error) {
		node, raw, err = cli.NodeInspectWithRaw(ctx, id)
		return
	})
	return
}

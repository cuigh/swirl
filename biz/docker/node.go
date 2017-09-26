package docker

import (
	"context"
	"sort"

	"github.com/cuigh/swirl/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// NodeList return all swarm nodes.
func NodeList() (infos []*model.NodeListInfo, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		var nodes []swarm.Node
		nodes, err = cli.NodeList(ctx, types.NodeListOptions{})
		if err == nil {
			sort.Slice(nodes, func(i, j int) bool {
				return nodes[i].Description.Hostname < nodes[j].Description.Hostname
			})
			infos = make([]*model.NodeListInfo, len(nodes))
			for i, n := range nodes {
				infos[i] = model.NewNodeListInfo(n)
			}
		}
		return
	})
	return
}

// NodeCount return number of swarm nodes.
func NodeCount() (count int, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		var nodes []swarm.Node
		nodes, err = cli.NodeList(ctx, types.NodeListOptions{})
		if err == nil {
			count = len(nodes)
		}
		return
	})
	return
}

// NodeRemove remove a swarm node from cluster.
func NodeRemove(id string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		return cli.NodeRemove(ctx, id, types.NodeRemoveOptions{})
	})
}

// NodeInspect return node information.
func NodeInspect(id string) (node swarm.Node, raw []byte, err error) {
	mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		node, raw, err = cli.NodeInspectWithRaw(ctx, id)
		return
	})
	return
}

// NodeUpdate update a node.
func NodeUpdate(id string, info *model.NodeUpdateInfo) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		version := swarm.Version{
			Index: info.Version,
		}
		spec := swarm.NodeSpec{
			Role:         info.Role,
			Availability: info.Availability,
		}
		spec.Name = info.Name
		spec.Labels = info.Labels.ToMap()
		return cli.NodeUpdate(ctx, id, version, spec)
	})
}

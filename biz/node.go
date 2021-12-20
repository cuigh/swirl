package biz

import (
	"context"
	"sort"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/docker"
	"github.com/docker/docker/api/types/swarm"
)

type NodeBiz interface {
	List() ([]*docker.Node, error)
	Search() ([]*Node, error)
	Find(id string) (node *Node, raw string, err error)
	Delete(id, name string, user web.User) (err error)
	Update(n *Node, user web.User) (err error)
}

func NewNode(d *docker.Docker, eb EventBiz) NodeBiz {
	return &nodeBiz{d: d, eb: eb}
}

type nodeBiz struct {
	d  *docker.Docker
	eb EventBiz
}

func (b *nodeBiz) Find(id string) (node *Node, raw string, err error) {
	var (
		sn swarm.Node
		r  []byte
	)
	sn, r, err = b.d.NodeInspect(context.TODO(), id)
	if err == nil {
		raw, err = indentJSON(r)
	}
	if err == nil {
		node = newNode(&sn)
	}
	return
}

func (b *nodeBiz) List() ([]*docker.Node, error) {
	m, err := b.d.NodeMap()
	if err != nil {
		return nil, err
	}

	nodes := make([]*docker.Node, 0, len(m))
	for _, n := range m {
		nodes = append(nodes, n)
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})
	return nodes, nil
}

func (b *nodeBiz) Search() ([]*Node, error) {
	list, err := b.d.NodeList(context.TODO())
	if err != nil {
		return nil, err
	}

	networks := make([]*Node, len(list))
	for i, n := range list {
		networks[i] = newNode(&n)
	}
	return networks, nil
}

func (b *nodeBiz) Delete(id, name string, user web.User) (err error) {
	err = b.d.NodeRemove(context.TODO(), id)
	if err == nil {
		b.eb.CreateNode(EventActionDelete, id, name, user)
	}
	return
}

func (b *nodeBiz) Update(n *Node, user web.User) (err error) {
	spec := &swarm.NodeSpec{
		Role:         n.Role,
		Availability: n.Availability,
	}
	spec.Name = n.Name
	spec.Labels = toMap(n.Labels)
	err = b.d.NodeUpdate(context.TODO(), n.ID, n.Version, spec)
	if err == nil {
		b.eb.CreateNode(EventActionUpdate, n.ID, n.Hostname, user)
	}
	return
}

type Node struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name,omitempty"`
	Hostname      string                 `json:"hostname"`
	Version       uint64                 `json:"version"`
	Role          swarm.NodeRole         `json:"role"`
	Availability  swarm.NodeAvailability `json:"availability"`
	EngineVersion string                 `json:"engineVersion"`
	Architecture  string                 `json:"arch"`
	OS            string                 `json:"os"`
	CPU           int64                  `json:"cpu"`
	Memory        float32                `json:"memory"`
	Address       string                 `json:"address"`
	State         swarm.NodeState        `json:"state"`
	Manager       *NodeManager           `json:"manager,omitempty"`
	Labels        data.Options           `json:"labels,omitempty"`
	CreatedAt     string                 `json:"createdAt"`
	UpdatedAt     string                 `json:"updatedAt"`
}

type NodeManager struct {
	Leader       bool               `json:"leader,omitempty"`
	Reachability swarm.Reachability `json:"reachability,omitempty"`
	Addr         string             `json:"addr,omitempty"`
}

func newNode(sn *swarm.Node) *Node {
	n := &Node{
		ID:            sn.ID,
		Name:          sn.Spec.Name,
		Hostname:      sn.Description.Hostname,
		Version:       sn.Version.Index,
		Role:          sn.Spec.Role,
		Availability:  sn.Spec.Availability,
		EngineVersion: sn.Description.Engine.EngineVersion,
		Architecture:  sn.Description.Platform.Architecture,
		OS:            sn.Description.Platform.OS,
		CPU:           sn.Description.Resources.NanoCPUs / 1e9,
		Memory:        float32(sn.Description.Resources.MemoryBytes>>20) / 1024,
		Address:       sn.Status.Addr,
		State:         sn.Status.State,
		Labels:        mapToOptions(sn.Spec.Labels),
		CreatedAt:     formatTime(sn.CreatedAt),
		UpdatedAt:     formatTime(sn.UpdatedAt),
	}
	if n.Role == swarm.NodeRoleManager {
		n.Manager = &NodeManager{
			Leader:       sn.ManagerStatus.Leader,
			Reachability: sn.ManagerStatus.Reachability,
			Addr:         sn.ManagerStatus.Addr,
		}
	}
	return n
}

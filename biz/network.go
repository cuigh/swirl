package biz

import (
	"context"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

type NetworkBiz interface {
	Search() ([]*Network, error)
	Find(name string) (network *Network, raw string, err error)
	Delete(id, name string, user web.User) (err error)
	Create(n *Network, user web.User) (err error)
	Disconnect(networkId, networkName, container string, user web.User) (err error)
}

func NewNetwork(d *docker.Docker, eb EventBiz) NetworkBiz {
	return &networkBiz{d: d, eb: eb}
}

type networkBiz struct {
	d  *docker.Docker
	eb EventBiz
}

func (b *networkBiz) Create(n *Network, user web.User) (err error) {
	nc := &types.NetworkCreate{
		Driver:     n.Driver,
		Scope:      n.Scope,
		Internal:   n.Internal,
		Attachable: n.Attachable,
		Ingress:    n.Ingress,
		EnableIPv6: n.IPv6,
		IPAM:       &network.IPAM{},
		Options:    toMap(n.Options),
		Labels:     toMap(n.Labels),
		//ConfigOnly     bool
		//ConfigFrom     *network.ConfigReference
	}
	for _, c := range n.IPAM.Config {
		nc.IPAM.Config = append(nc.IPAM.Config, network.IPAMConfig{
			Subnet:  c.Subnet,
			Gateway: c.Gateway,
			IPRange: c.Range,
		})
	}
	err = b.d.NetworkCreate(context.TODO(), n.Name, nc)
	if err != nil {
		b.eb.CreateNetwork(EventActionCreate, n.Name, n.Name, user)
	}
	return
}

func (b *networkBiz) Find(name string) (network *Network, raw string, err error) {
	var (
		nr types.NetworkResource
		r  []byte
	)
	nr, r, err = b.d.NetworkInspect(context.TODO(), name)
	if err == nil {
		network = newNetwork(&nr)
		raw, err = indentJSON(r)
	}
	return
}

func (b *networkBiz) Search() ([]*Network, error) {
	list, err := b.d.NetworkList(context.TODO())
	if err != nil {
		return nil, err
	}

	networks := make([]*Network, len(list))
	for i, nr := range list {
		networks[i] = newNetwork(&nr)
	}
	return networks, nil
}

func (b *networkBiz) Delete(id, name string, user web.User) (err error) {
	err = b.d.NetworkRemove(context.TODO(), name)
	if err == nil {
		b.eb.CreateNetwork(EventActionDelete, id, name, user)
	}
	return
}

func (b *networkBiz) Disconnect(networkId, networkName, container string, user web.User) (err error) {
	err = b.d.NetworkDisconnect(context.TODO(), networkName, container)
	if err == nil {
		b.eb.CreateNetwork(EventActionDisconnect, networkId, networkName, user)
	}
	return
}

type Network struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Created    string `json:"created"`
	Driver     string `json:"driver"`
	Scope      string `json:"scope"`
	Internal   bool   `json:"internal"`
	Attachable bool   `json:"attachable"`
	Ingress    bool   `json:"ingress"`
	IPv6       bool   `json:"ipv6"`
	IPAM       struct {
		Driver  string        `json:"driver"`
		Options data.Options  `json:"options"`
		Config  []*IPAMConfig `json:"config"`
	} `json:"ipam"`
	Options    data.Options        `json:"options"`
	Labels     data.Options        `json:"labels"`
	Containers []*NetworkContainer `json:"containers"`
}

type IPAMConfig struct {
	Subnet  string `json:"subnet,omitempty"`
	Gateway string `json:"gateway,omitempty"`
	Range   string `json:"range,omitempty"`
}

type NetworkContainer struct {
	ID   string `json:"id"`   // container id
	Name string `json:"name"` // container name
	Mac  string `json:"mac"`  // mac address
	IPv4 string `json:"ipv4"` // IPv4 address
	IPv6 string `json:"ipv6"` // IPv6 address
}

func newNetwork(nr *types.NetworkResource) *Network {
	n := &Network{
		ID:         nr.ID,
		Name:       nr.Name,
		Created:    formatTime(nr.Created),
		Driver:     nr.Driver,
		Scope:      nr.Scope,
		Internal:   nr.Internal,
		Attachable: nr.Attachable,
		Ingress:    nr.Ingress,
		IPv6:       nr.EnableIPv6,
		Options:    mapToOptions(nr.Options),
		Labels:     mapToOptions(nr.Labels),
	}
	n.IPAM.Driver = nr.IPAM.Driver
	n.IPAM.Options = mapToOptions(nr.IPAM.Options)
	n.IPAM.Config = make([]*IPAMConfig, len(nr.IPAM.Config))
	for i, c := range nr.IPAM.Config {
		n.IPAM.Config[i] = &IPAMConfig{
			Subnet:  c.Subnet,
			Gateway: c.Gateway,
			Range:   c.IPRange,
		}
	}
	n.Containers = make([]*NetworkContainer, 0, len(nr.Containers))
	for id, ep := range nr.Containers {
		n.Containers = append(n.Containers, &NetworkContainer{
			ID:   id,
			Name: ep.Name,
			Mac:  ep.MacAddress,
			IPv4: ep.IPv4Address,
			IPv6: ep.IPv6Address,
		})
	}
	return n
}

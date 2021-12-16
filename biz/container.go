package biz

import (
	"context"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
)

type Container struct {
	/* Summary */
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Image       string            `json:"image,omitempty"`
	Command     string            `json:"command,omitempty"`
	CreatedAt   string            `json:"createdAt"`
	Ports       []*ContainerPort  `json:"ports,omitempty"`
	SizeRw      int64             `json:"sizeRw"`
	SizeRootFs  int64             `json:"sizeRootFs"`
	Labels      data.Options      `json:"labels"`
	State       string            `json:"state"`
	Status      string            `json:"status"`
	NetworkMode string            `json:"networkMode"`
	Mounts      []*ContainerMount `json:"mounts"`
	//Networks map[string]*network.EndpointSettings

	/* Detail */
	PID       int    `json:"pid,omitempty"`
	StartedAt string `json:"startedAt,omitempty"`
}

type ContainerPort struct {
	IP          string `json:"ip,omitempty"`
	PrivatePort uint16 `json:"privatePort,omitempty"`
	PublicPort  uint16 `json:"publicPort,omitempty"`
	Type        string `json:"type,omitempty"`
}

type ContainerMount struct {
	Type        mount.Type        `json:"type,omitempty"`
	Name        string            `json:"name,omitempty"`
	Source      string            `json:"source,omitempty"`
	Destination string            `json:"destination,omitempty"`
	Driver      string            `json:"driver,omitempty"`
	Mode        string            `json:"mode,omitempty"`
	RW          bool              `json:"rw,omitempty"`
	Propagation mount.Propagation `json:"propagation,omitempty"`
}

func newContainerMount(m types.MountPoint) *ContainerMount {
	return &ContainerMount{
		Type:        m.Type,
		Name:        m.Name,
		Source:      m.Source,
		Destination: m.Destination,
		Driver:      m.Driver,
		Mode:        m.Mode,
		RW:          m.RW,
		Propagation: m.Propagation,
	}
}

// EndpointSettings stores the network endpoint details
//type EndpointSettings struct {
//	// Configurations
//	IPAMConfig *EndpointIPAMConfig
//	Links      []string
//	Aliases    []string
//	// Operational data
//	NetworkID           string
//	EndpointID          string
//	Gateway             string
//	IPAddress           string
//	IPPrefixLen         int
//	IPv6Gateway         string
//	GlobalIPv6Address   string
//	GlobalIPv6PrefixLen int
//	MacAddress          string
//	DriverOpts          map[string]string
//}

func newContainerSummary(c *types.Container) *Container {
	container := &Container{
		ID:          c.ID,
		Name:        c.Names[0],
		Image:       normalizeImage(c.Image),
		Command:     c.Command,
		CreatedAt:   formatTime(time.Unix(c.Created, 0)),
		SizeRw:      c.SizeRw,
		SizeRootFs:  c.SizeRootFs,
		Labels:      mapToOptions(c.Labels),
		State:       c.State,
		Status:      c.Status,
		NetworkMode: c.HostConfig.NetworkMode,
	}
	for _, p := range c.Ports {
		container.Ports = append(container.Ports, &ContainerPort{
			IP:          p.IP,
			PrivatePort: p.PrivatePort,
			PublicPort:  p.PublicPort,
			Type:        p.Type,
		})
	}
	for _, m := range c.Mounts {
		container.Mounts = append(container.Mounts, newContainerMount(m))
	}
	return container
}

func newContainerDetail(c *types.ContainerJSON) *Container {
	created, _ := time.Parse(time.RFC3339Nano, c.Created)
	startedAt, _ := time.Parse(time.RFC3339Nano, c.State.StartedAt)
	container := &Container{
		ID:    c.ID,
		Name:  c.Name,
		Image: c.Image,
		//Command:     c.Command,
		CreatedAt: formatTime(created),
		Labels:    mapToOptions(c.Config.Labels),
		State:     c.State.Status,
		//Status:      c.Status,
		NetworkMode: string(c.HostConfig.NetworkMode),
		// detail
		PID:       c.State.Pid,
		StartedAt: formatTime(startedAt),
	}
	if c.SizeRw != nil {
		container.SizeRw = *c.SizeRw
	}
	if c.SizeRootFs != nil {
		container.SizeRootFs = *c.SizeRootFs
	}
	//for _, p := range c.Ports {
	//	container.Ports = append(container.Ports, &ContainerPort{
	//		IP:          p.IP,
	//		PrivatePort: p.PrivatePort,
	//		PublicPort:  p.PublicPort,
	//		Type:        p.Type,
	//	})
	//}
	for _, m := range c.Mounts {
		container.Mounts = append(container.Mounts, newContainerMount(m))
	}
	return container
}

type ContainerBiz interface {
	Search(name, status string, pageIndex, pageSize int) ([]*Container, int, error)
	Find(id string) (container *Container, raw string, err error)
	Delete(id string, user web.User) (err error)
	FetchLogs(id string, lines int, timestamps bool) (stdout, stderr string, err error)
	ExecCreate(id string, cmd string) (resp types.IDResponse, err error)
	ExecAttach(id string) (resp types.HijackedResponse, err error)
	ExecStart(id string) error
}

func NewContainer(d *docker.Docker) ContainerBiz {
	return &containerBiz{d: d}
}

type containerBiz struct {
	d *docker.Docker
}

func (b *containerBiz) Find(id string) (c *Container, raw string, err error) {
	var (
		cj types.ContainerJSON
		r  []byte
	)

	if cj, r, err = b.d.ContainerInspect(context.TODO(), id); err == nil {
		raw, err = indentJSON(r)
	}

	if err == nil {
		c = newContainerDetail(&cj)
	}
	return
}

func (b *containerBiz) Search(name, status string, pageIndex, pageSize int) (containers []*Container, total int, err error) {
	list, total, err := b.d.ContainerList(context.TODO(), name, status, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	containers = make([]*Container, len(list))
	for i, nr := range list {
		containers[i] = newContainerSummary(&nr)
	}
	return containers, total, nil
}

func (b *containerBiz) Delete(id string, user web.User) (err error) {
	err = b.d.ContainerRemove(context.TODO(), id)
	//if err == nil {
	//	Event.CreateContainer(model.EventActionDelete, id, user)
	//}
	return
}

func (b *containerBiz) ExecCreate(id, cmd string) (resp types.IDResponse, err error) {
	return b.d.ContainerExecCreate(context.TODO(), id, cmd)
}

func (b *containerBiz) ExecAttach(id string) (resp types.HijackedResponse, err error) {
	return b.d.ContainerExecAttach(context.TODO(), id)
}

func (b *containerBiz) ExecStart(id string) error {
	return b.d.ContainerExecStart(context.TODO(), id)
}

func (b *containerBiz) FetchLogs(id string, lines int, timestamps bool) (string, string, error) {
	stdout, stderr, err := b.d.ContainerLogs(context.TODO(), id, lines, timestamps)
	if err != nil {
		return "", "", err
	}
	return stdout.String(), stderr.String(), nil
}
package biz

import (
	"context"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/swirl/docker"
	"github.com/docker/docker/api/types/swarm"
)

type TaskBiz interface {
	Search(node, service, mode string, pageIndex, pageSize int) (tasks []*Task, total int, err error)
	Find(id string) (task *Task, raw string, err error)
	FetchLogs(id string, lines int, timestamps bool) (stdout, stderr string, err error)
}

func NewTask(d *docker.Docker) TaskBiz {
	return &taskBiz{d: d}
}

type taskBiz struct {
	d *docker.Docker
}

func (b *taskBiz) Find(id string) (task *Task, raw string, err error) {
	var (
		t swarm.Task
		s swarm.Service
		r []byte
	)

	t, r, err = b.d.TaskInspect(context.TODO(), id)
	if err == nil {
		raw, err = indentJSON(r)
	}

	if err == nil {
		m, _ := b.d.NodeMap()
		task = newTask(&t, m)

		// Fill service name
		if s, _, _ = b.d.ServiceInspect(context.TODO(), t.ServiceID, false); s.Spec.Name == "" {
			task.ServiceName = task.ServiceID
		} else {
			task.ServiceName = s.Spec.Name
		}
	}
	return
}

func (b *taskBiz) Search(node, service, state string, pageIndex, pageSize int) (tasks []*Task, total int, err error) {
	var list []swarm.Task
	list, total, err = b.d.TaskList(context.TODO(), node, service, state, pageIndex, pageSize)
	if err != nil {
		return
	}

	m, _ := b.d.NodeMap()
	tasks = make([]*Task, len(list))
	for i, t := range list {
		tasks[i] = newTask(&t, m)
		if m != nil {
			if n, ok := m[t.NodeID]; ok {
				tasks[i].NodeName = n.Name
			}
		}
	}
	return
}

func (b *taskBiz) FetchLogs(id string, lines int, timestamps bool) (string, string, error) {
	stdout, stderr, err := b.d.TaskLogs(context.TODO(), id, lines, timestamps)
	if err != nil {
		return "", "", err
	}
	return stdout.String(), stderr.String(), nil
}

type Task struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Version     uint64          `json:"version"`
	Image       string          `json:"image"`
	Slot        int             `json:"slot"`
	State       swarm.TaskState `json:"state"`
	ServiceID   string          `json:"serviceId"`
	ServiceName string          `json:"serviceName"`
	NodeID      string          `json:"nodeId"`
	NodeName    string          `json:"nodeName"`
	ContainerID string          `json:"containerId"`
	PID         int             `json:"pid"`
	ExitCode    int             `json:"exitCode"`
	Message     string          `json:"message"`
	Error       string          `json:"error"`
	Env         data.Options    `json:"env,omitempty"`
	Labels      data.Options    `json:"labels,omitempty"`
	Networks    []TaskNetwork   `json:"networks"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
}

type TaskNetwork struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	IPs  []string `json:"ips"`
}

func newTask(t *swarm.Task, nodes map[string]*docker.Node) *Task {
	task := &Task{
		ID:        t.ID,
		Name:      t.Name,
		Version:   t.Version.Index,
		Image:     normalizeImage(t.Spec.ContainerSpec.Image),
		Slot:      t.Slot,
		State:     t.Status.State,
		ServiceID: t.ServiceID,
		NodeID:    t.NodeID,
		NodeName:  t.NodeID,
		Message:   t.Status.Message,
		Error:     t.Status.Err,
		Env:       envToOptions(t.Spec.ContainerSpec.Env),
		Labels:    mapToOptions(t.Labels),
		CreatedAt: formatTime(t.CreatedAt),
		UpdatedAt: formatTime(t.UpdatedAt),
	}
	if t.Status.ContainerStatus != nil {
		task.ContainerID = t.Status.ContainerStatus.ContainerID
		task.PID = t.Status.ContainerStatus.PID
		task.ExitCode = t.Status.ContainerStatus.ExitCode
	}
	for _, n := range t.NetworksAttachments {
		task.Networks = append(task.Networks, TaskNetwork{
			ID:   n.Network.ID,
			Name: n.Network.Spec.Name,
			IPs:  n.Addresses,
		})
	}
	// Fill node name
	if nodes != nil {
		if n, ok := nodes[t.NodeID]; ok {
			task.NodeName = n.Name
		}
	}
	return task
}

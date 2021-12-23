package docker

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/cache"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/util/lazy"
	"github.com/cuigh/swirl/misc"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

func newVersion(v uint64) swarm.Version {
	return swarm.Version{Index: v}
}

type Docker struct {
	c        *client.Client
	locker   sync.Mutex
	logger   log.Logger
	nodes    cache.Value
	agents   sync.Map
	networks sync.Map
}

func NewDocker() *Docker {
	d := &Docker{
		logger: log.Get("docker"),
		nodes:  cache.Value{TTL: 30 * time.Minute},
	}
	d.nodes.Load = d.loadCache
	return d
}

func (d *Docker) call(fn func(c *client.Client) error) error {
	c, err := d.client()
	if err == nil {
		err = fn(c)
	}
	return err
}

func (d *Docker) client() (c *client.Client, err error) {
	if d.c == nil {
		d.locker.Lock()
		defer d.locker.Unlock()

		if d.c == nil {
			var opt client.Opt
			if misc.Options.DockerEndpoint == "" {
				opt = client.FromEnv
			} else {
				opt = client.WithHost(misc.Options.DockerEndpoint)
			}
			d.c, err = client.NewClientWithOpts(opt, client.WithVersion(misc.Options.DockerAPIVersion))
			if err != nil {
				return
			}
		}
	}
	return d.c, nil
}

func (d *Docker) agent(node string) (*client.Client, error) {
	host, err := d.getAgent(node)
	if err != nil {
		d.logger.Error("failed to find node agent: ", err)
	}

	if host == "" {
		return d.client()
	}

	value, _ := d.agents.LoadOrStore(node, &lazy.Value{
		New: func() (interface{}, error) {
			c, e := client.NewClientWithOpts(
				client.WithHost("tcp://"+host),
				client.WithVersion(misc.Options.DockerAPIVersion),
			)
			return c, e
		},
	})
	c, err := value.(*lazy.Value).Get()
	if err != nil {
		return nil, err
	}
	return c.(*client.Client), nil
}

func (d *Docker) getAgent(node string) (agent string, err error) {
	if node == "" || node == "-" {
		return "", nil
	}

	nodes, err := d.NodeMap()
	if err != nil {
		return
	}

	if n, ok := nodes[node]; ok {
		agent = n.Agent
	}
	return
}

func (d *Docker) loadCache() (interface{}, error) {
	c, err := d.client()
	if err != nil {
		return nil, err
	}

	agents, err := d.loadAgents(context.TODO(), c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load agents")
	}

	nodes, err := d.loadNodes(context.TODO(), c)
	if err != nil {
		return nil, err
	}
	for i := range nodes {
		nodes[i].Agent = agents[nodes[i].ID]
	}

	return nodes, nil
}

func (d *Docker) loadNodes(ctx context.Context, c *client.Client) (nodes map[string]*Node, err error) {
	var list []swarm.Node
	list, err = c.NodeList(ctx, types.NodeListOptions{})
	if err == nil {
		nodes = make(map[string]*Node)
		for _, n := range list {
			ni := &Node{
				ID:    n.ID,
				Name:  n.Spec.Name,
				State: n.Status.State,
			}
			if ni.Name == "" {
				ni.Name = n.Description.Hostname
			}
			nodes[n.ID] = ni
		}
	}
	return
}

func (d *Docker) loadAgents(ctx context.Context, c *client.Client) (agents map[string]string, err error) {
	var tasks []swarm.Task
	agents = make(map[string]string)
	for _, agent := range misc.Options.Agents {
		pair := strings.SplitN(agent, ":", 2)
		args := filters.NewArgs(
			filters.Arg("desired-state", string(swarm.TaskStateRunning)),
			filters.Arg("service", pair[0]),
		)
		tasks, err = c.TaskList(ctx, types.TaskListOptions{Filters: args})
		if err != nil {
			return
		}

		port := "2375"
		if len(pair) > 1 {
			port = pair[1]
		}

		for _, t := range tasks {
			if len(t.NetworksAttachments) > 0 {
				pair = strings.SplitN(t.NetworksAttachments[0].Addresses[0], "/", 2)
				agents[t.NodeID] = pair[0] + ":" + port
			}
		}
	}
	return
}

type Node struct {
	ID    string          `json:"id,omitempty"`
	Name  string          `json:"name,omitempty"`
	State swarm.NodeState `json:"-"`
	Agent string          `json:"-"`
}

func init() {
	container.Put(NewDocker)
}

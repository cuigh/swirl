package docker

import (
	"os"
	"sync"

	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/swirl/misc"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

const (
	defaultAPIVersion = "1.41"
)

func newVersion(v uint64) swarm.Version {
	return swarm.Version{Index: v}
}

type Docker struct {
	c      *client.Client
	locker sync.Mutex
	logger log.Logger
}

func NewDocker() *Docker {
	return &Docker{
		logger: log.Get("docker"),
	}
}

func (d *Docker) call(fn func(c *client.Client) error) error {
	c, err := d.client()
	if err == nil {
		err = fn(c)
	}
	return err
}

func (d *Docker) client() (cli *client.Client, err error) {
	if d.c == nil {
		d.locker.Lock()
		defer d.locker.Unlock()

		if d.c == nil {
			apiVersion := misc.Options.DockerAPIVersion
			if apiVersion == "" {
				apiVersion = defaultAPIVersion
			}
			if misc.Options.DockerEndpoint == "" {
				_ = os.Setenv("DOCKER_API_VERSION", apiVersion)
				d.c, err = client.NewClientWithOpts(client.FromEnv)
			} else {
				d.c, err = client.NewClientWithOpts(client.WithHost(misc.Options.DockerEndpoint), client.WithVersion(apiVersion))
			}
			if err != nil {
				return
			}
		}
	}
	return d.c, nil
}

func init() {
	container.Put(NewDocker)
}

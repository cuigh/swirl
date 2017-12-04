package docker

import (
	"context"
	"os"
	"sync"

	"github.com/cuigh/auxo/log"
	"github.com/cuigh/swirl/misc"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

const (
	apiVersion = "1.32"
)

var mgr = &manager{}

type manager struct {
	client *client.Client
	locker sync.Mutex
	logger *log.Logger
}

func (m *manager) Do(fn func(ctx context.Context, cli *client.Client) error) (err error) {
	ctx, cli, err := m.Client()
	if err != nil {
		return err
	}
	return fn(ctx, cli)
}

func (m *manager) Client() (ctx context.Context, cli *client.Client, err error) {
	if m.client == nil {
		m.locker.Lock()
		defer m.locker.Unlock()

		if m.client == nil {
			if misc.Options.DockerEndpoint == "" {
				os.Setenv("DOCKER_API_VERSION", apiVersion)
				m.client, err = client.NewEnvClient()
			} else {
				m.client, err = client.NewClient(misc.Options.DockerEndpoint, apiVersion, nil, nil)
			}
			if err != nil {
				return
			}
		}
	}
	return context.TODO(), m.client, nil
}

func (m *manager) Logger() *log.Logger {
	if m.logger == nil {
		m.locker.Lock()
		defer m.locker.Unlock()

		if m.logger == nil {
			m.logger = log.Get("docker")
		}
	}
	return m.logger
}

func version(v uint64) swarm.Version {
	return swarm.Version{Index: v}
}

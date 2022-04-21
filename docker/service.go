package docker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sort"
	"strconv"

	"github.com/cuigh/swirl/misc"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// ServiceList return service list.
func (d *Docker) ServiceList(ctx context.Context, name, mode string, pageIndex, pageSize int) (services []swarm.Service, totalCount int, err error) {
	err = d.call(func(c *client.Client) (err error) {
		args := filters.NewArgs()
		if name != "" {
			args.Add("name", name)
		}
		if mode != "" {
			args.Add("mode", mode)
		}
		services, err = c.ServiceList(ctx, types.ServiceListOptions{Status: true, Filters: args})
		if err != nil {
			return
		}

		sort.Slice(services, func(i, j int) bool {
			return services[i].Spec.Name < services[j].Spec.Name
		})
		totalCount = len(services)
		start, end := misc.Page(totalCount, pageIndex, pageSize)
		services = services[start:end]

		if versions.LessThan(c.ClientVersion(), "1.41") {
			err = d.fillStatus(ctx, c, services)
		}
		return
	})
	return
}

// ServiceInspect return service raw information.
func (d *Docker) ServiceInspect(ctx context.Context, name string, status bool) (service swarm.Service, raw []byte, err error) {
	var c *client.Client
	if c, err = d.client(); err != nil {
		return
	}

	service, raw, err = c.ServiceInspectWithRaw(ctx, name, types.ServiceInspectOptions{})
	if err == nil && status {
		var services []swarm.Service
		if versions.LessThan(c.ClientVersion(), "1.41") {
			services = []swarm.Service{service}
			_ = d.fillStatus(ctx, c, services)
			service.ServiceStatus = services[0].ServiceStatus
		} else {
			services, _ = c.ServiceList(ctx, types.ServiceListOptions{
				Status:  true,
				Filters: filters.NewArgs(filters.Arg("name", name)),
			})
			if len(services) > 0 {
				service.ServiceStatus = services[0].ServiceStatus
			}
		}
	}
	return
}

func (d *Docker) fillStatus(ctx context.Context, c *client.Client, services []swarm.Service) (err error) {
	var (
		nodes map[string]*Node
		m     = make(map[string]*swarm.Service)
		tasks []swarm.Task
		opts  = types.TaskListOptions{Filters: filters.NewArgs()}
	)

	nodes, err = d.NodeMap()
	if err != nil {
		return
	}

	for i := range services {
		s := &services[i]
		s.ServiceStatus = &swarm.ServiceStatus{}
		if s.Spec.Mode.Replicated != nil {
			s.ServiceStatus.DesiredTasks = *s.Spec.Mode.Replicated.Replicas
		}
		m[s.ID] = s
		opts.Filters.Add("service", s.ID)
	}

	// acquire all related tasks
	tasks, err = c.TaskList(ctx, opts)
	if err != nil {
		return
	}

	// count tasks
	for _, task := range tasks {
		s := m[task.ServiceID]
		if s != nil {
			if s.Spec.Mode.Global != nil && task.DesiredState != swarm.TaskStateShutdown {
				s.ServiceStatus.DesiredTasks++
			}
			if n, ok := nodes[task.NodeID]; ok && n.State != swarm.NodeStateDown && task.Status.State == swarm.TaskStateRunning {
				s.ServiceStatus.RunningTasks++
			}
		}
	}
	return
}

// ServiceRemove remove a service.
func (d *Docker) ServiceRemove(ctx context.Context, name string) error {
	return d.call(func(c *client.Client) (err error) {
		return c.ServiceRemove(ctx, name)
	})
}

// ServiceRollback rollbacks a service.
func (d *Docker) ServiceRollback(ctx context.Context, name string) error {
	return d.call(func(c *client.Client) (err error) {
		service, _, err := c.ServiceInspectWithRaw(ctx, name, types.ServiceInspectOptions{})
		if err != nil {
			return err
		}

		options := types.ServiceUpdateOptions{
			Rollback: "previous",
		}
		resp, err := c.ServiceUpdate(ctx, name, service.Version, service.Spec, options)
		if err == nil && len(resp.Warnings) > 0 {
			d.logger.Warnf("service '%s' was rollbacked but got warnings: %v", name, resp.Warnings)
		}
		return err
	})
}

// ServiceRestart force to refresh a service.
func (d *Docker) ServiceRestart(ctx context.Context, name string) error {
	return d.call(func(c *client.Client) (err error) {
		service, _, err := c.ServiceInspectWithRaw(ctx, name, types.ServiceInspectOptions{})
		if err != nil {
			return err
		}

		service.Spec.TaskTemplate.ForceUpdate++
		resp, err := c.ServiceUpdate(ctx, name, service.Version, service.Spec, types.ServiceUpdateOptions{})
		if err == nil && len(resp.Warnings) > 0 {
			d.logger.Warnf("service '%s' was restarted but got warnings: %v", name, resp.Warnings)
		}
		return err
	})
}

// ServiceScale adjust replicas of a service.
func (d *Docker) ServiceScale(ctx context.Context, name string, count, version uint64) error {
	return d.call(func(c *client.Client) (err error) {
		service, _, err := c.ServiceInspectWithRaw(ctx, name, types.ServiceInspectOptions{})
		if err != nil {
			return err
		}

		spec := service.Spec
		if spec.Mode.Replicated != nil {
			spec.Mode.Replicated.Replicas = &count
		} else if spec.Mode.ReplicatedJob != nil {
			spec.Mode.ReplicatedJob.TotalCompletions = &count
		} else {
			return errors.New("scale can only be used with replicated or replicated-job mode")
		}

		ver := service.Version
		if version > 0 {
			ver = swarm.Version{Index: version}
		}
		resp, err := c.ServiceUpdate(ctx, name, ver, spec, types.ServiceUpdateOptions{})
		if err == nil && len(resp.Warnings) > 0 {
			d.logger.Warnf("service %s was scaled but got warnings: %v", name, resp.Warnings)
		}
		return err
	})
}

// ServiceCreate create a service.
func (d *Docker) ServiceCreate(ctx context.Context, spec *swarm.ServiceSpec, registryAuth string) error {
	return d.call(func(c *client.Client) (err error) {
		opts := types.ServiceCreateOptions{EncodedRegistryAuth: registryAuth}
		var resp types.ServiceCreateResponse
		resp, err = c.ServiceCreate(ctx, *spec, opts)
		if err == nil && len(resp.Warnings) > 0 {
			d.logger.Warnf("service %s was created but got warnings: %v", spec.Name, resp.Warnings)
		}
		return
	})
}

// ServiceUpdate update a service.
func (d *Docker) ServiceUpdate(ctx context.Context, spec *swarm.ServiceSpec, version uint64) error {
	return d.call(func(c *client.Client) (err error) {
		var (
			resp    types.ServiceUpdateResponse
			options = types.ServiceUpdateOptions{
				RegistryAuthFrom: types.RegistryAuthFromSpec,
				QueryRegistry:    false,
			}
		)
		resp, err = c.ServiceUpdate(ctx, spec.Name, newVersion(version), *spec, options)
		if err == nil && len(resp.Warnings) > 0 {
			d.logger.Warnf("service %s was updated but got warnings: %v", spec.Name, resp.Warnings)
		}
		return
	})
}

// ServiceCount return number of services.
func (d *Docker) ServiceCount(ctx context.Context) (count int, err error) {
	err = d.call(func(c *client.Client) (err error) {
		var services []swarm.Service
		if services, err = c.ServiceList(ctx, types.ServiceListOptions{}); err == nil {
			count = len(services)
		}
		return
	})
	return
}

// ServiceLogs returns the logs generated by a service.
func (d *Docker) ServiceLogs(ctx context.Context, name string, lines int, timestamps bool) (stdout, stderr *bytes.Buffer, err error) {
	err = d.call(func(c *client.Client) (err error) {
		var (
			rc   io.ReadCloser
			opts = types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Tail:       strconv.Itoa(lines),
				Timestamps: timestamps,
				//Since: (time.Hour * 24).String()
			}
		)
		if rc, err = c.ServiceLogs(ctx, name, opts); err == nil {
			defer rc.Close()

			stdout = &bytes.Buffer{}
			stderr = &bytes.Buffer{}
			_, err = stdcopy.StdCopy(stdout, stderr, rc)
		}
		return
	})
	return
}

// ServiceSearch search services with args.
func (d *Docker) ServiceSearch(ctx context.Context, args filters.Args) (services []swarm.Service, err error) {
	err = d.call(func(c *client.Client) (err error) {
		opts := types.ServiceListOptions{Filters: args}
		services, err = c.ServiceList(ctx, opts)
		return
	})
	return
}

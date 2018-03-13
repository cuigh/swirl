package docker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cuigh/auxo/ext/texts"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// ServiceList return service list.
func ServiceList(name string, pageIndex, pageSize int) (infos []*model.ServiceListInfo, totalCount int, err error) { // nolint: gocyclo
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		var (
			services    []swarm.Service
			nodes       []swarm.Node
			activeNodes map[string]struct{}
			tasks       []swarm.Task
		)

		// acquire services
		opts := types.ServiceListOptions{}
		if name != "" {
			opts.Filters = filters.NewArgs()
			opts.Filters.Add("name", name)
		}
		services, err = cli.ServiceList(ctx, opts)
		if err != nil {
			return
		}
		sort.Slice(services, func(i, j int) bool {
			return services[i].Spec.Name < services[j].Spec.Name
		})
		totalCount = len(services)
		start, end := misc.Page(totalCount, pageIndex, pageSize)
		services = services[start:end]

		// acquire all swarm nodes
		nodes, err = cli.NodeList(ctx, types.NodeListOptions{})
		if err != nil {
			return
		}
		activeNodes = make(map[string]struct{})
		for _, n := range nodes {
			if n.Status.State != swarm.NodeStateDown {
				activeNodes[n.ID] = struct{}{}
			}
		}

		// acquire all related tasks
		taskOpts := types.TaskListOptions{
			Filters: filters.NewArgs(),
		}
		for _, service := range services {
			taskOpts.Filters.Add("service", service.ID)
		}
		tasks, err = cli.TaskList(ctx, taskOpts)
		if err != nil {
			return
		}

		// count active tasks
		running, tasksNoShutdown := map[string]uint64{}, map[string]uint64{}
		for _, task := range tasks {
			if task.DesiredState != swarm.TaskStateShutdown {
				tasksNoShutdown[task.ServiceID]++
			}

			if _, nodeActive := activeNodes[task.NodeID]; nodeActive && task.Status.State == swarm.TaskStateRunning {
				running[task.ServiceID]++
			}
		}

		infos = make([]*model.ServiceListInfo, len(services))
		for i, service := range services {
			infos[i] = model.NewServiceListInfo(service, running[service.ID])
			if service.Spec.Mode.Global != nil {
				infos[i].Replicas = tasksNoShutdown[service.ID]
			}
		}
		return
	})
	return
}

// ServiceSearch search services with args.
func ServiceSearch(args filters.Args) (services []swarm.Service, err error) { // nolint: gocyclo
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		opts := types.ServiceListOptions{Filters: args}
		services, err = cli.ServiceList(ctx, opts)
		return
	})
	return
}

// ServiceCount return number of services.
func ServiceCount() (count int, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		var services []swarm.Service
		if services, err = cli.ServiceList(ctx, types.ServiceListOptions{}); err == nil {
			count = len(services)
		}
		return
	})
	return
}

// ServiceInspect return service raw information.
func ServiceInspect(name string) (service swarm.Service, raw []byte, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		service, raw, err = cli.ServiceInspectWithRaw(ctx, name, types.ServiceInspectOptions{})
		return
	})
	return
}

// ServiceUpdate update a service.
func ServiceUpdate(info *model.ServiceInfo) error { // nolint: gocyclo
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		service, _, err := cli.ServiceInspectWithRaw(ctx, info.Name, types.ServiceInspectOptions{})
		if err != nil {
			return err
		}
		spec := service.Spec

		// Annotations
		spec.Annotations.Labels = info.ServiceLabels.ToMap()

		// ContainerSpec
		spec.TaskTemplate.ContainerSpec.Image = info.Image
		spec.TaskTemplate.ContainerSpec.Dir = info.Dir
		spec.TaskTemplate.ContainerSpec.User = info.User
		spec.TaskTemplate.ContainerSpec.Labels = info.ContainerLabels.ToMap()
		spec.TaskTemplate.ContainerSpec.Command = nil
		spec.TaskTemplate.ContainerSpec.Args = nil
		spec.TaskTemplate.ContainerSpec.Env = nil
		if info.Command != "" {
			spec.TaskTemplate.ContainerSpec.Command = strings.Split(info.Command, " ")
		}
		if info.Args != "" {
			spec.TaskTemplate.ContainerSpec.Args = strings.Split(info.Args, " ")
		}
		if envs := info.Environments.ToMap(); len(envs) > 0 {
			for n, v := range envs {
				spec.TaskTemplate.ContainerSpec.Env = append(spec.TaskTemplate.ContainerSpec.Env, n+"="+v)
			}
		}

		// Mode
		if info.Mode == "replicated" {
			if spec.Mode.Replicated == nil {
				spec.Mode.Global = nil
				spec.Mode.Replicated = &swarm.ReplicatedService{Replicas: &info.Replicas}
			} else {
				spec.Mode.Replicated.Replicas = &info.Replicas
			}
		} else if info.Mode == "global" && spec.Mode.Global == nil {
			spec.Mode.Replicated = nil
			spec.Mode.Global = &swarm.GlobalService{}
		}

		// Network
		// todo: only process updated networks
		spec.TaskTemplate.Networks = nil
		for _, n := range info.Networks {
			spec.TaskTemplate.Networks = append(spec.TaskTemplate.Networks, swarm.NetworkAttachmentConfig{Target: n})
		}

		// Endpoint
		spec.EndpointSpec = &swarm.EndpointSpec{Mode: swarm.ResolutionMode(info.Endpoint.Mode)}
		for _, p := range info.Endpoint.Ports {
			port := swarm.PortConfig{
				Protocol:      p.Protocol,
				TargetPort:    p.TargetPort,
				PublishedPort: p.PublishedPort,
				PublishMode:   p.PublishMode,
			}
			spec.EndpointSpec.Ports = append(spec.EndpointSpec.Ports, port)
		}

		spec.TaskTemplate.ContainerSpec.Secrets = nil
		for _, s := range info.Secrets {
			spec.TaskTemplate.ContainerSpec.Secrets = append(spec.TaskTemplate.ContainerSpec.Secrets, s.ToSecret())
		}

		spec.TaskTemplate.ContainerSpec.Configs = nil
		for _, c := range info.Configs {
			spec.TaskTemplate.ContainerSpec.Configs = append(spec.TaskTemplate.ContainerSpec.Configs, c.ToConfig())
		}

		// Mounts
		// todo: fix > original options are not reserved.
		spec.TaskTemplate.ContainerSpec.Mounts = nil
		for _, m := range info.Mounts {
			if m.Target != "" {
				mnt := mount.Mount{
					Type:     mount.Type(m.Type),
					Source:   m.Source,
					Target:   m.Target,
					ReadOnly: m.ReadOnly,
				}
				if m.Propagation != "" {
					mnt.BindOptions = &mount.BindOptions{Propagation: m.Propagation}
				}
				spec.TaskTemplate.ContainerSpec.Mounts = append(spec.TaskTemplate.ContainerSpec.Mounts, mnt)
			}
		}

		// Placement
		if spec.TaskTemplate.Placement == nil {
			spec.TaskTemplate.Placement = &swarm.Placement{}
		}
		spec.TaskTemplate.Placement.Constraints = nil
		for _, c := range info.Placement.Constraints {
			if cons := c.ToConstraint(); cons != "" {
				spec.TaskTemplate.Placement.Constraints = append(spec.TaskTemplate.Placement.Constraints, cons)
			}
		}
		spec.TaskTemplate.Placement.Preferences = nil
		for _, p := range info.Placement.Preferences {
			if p.Spread != "" {
				pref := swarm.PlacementPreference{
					Spread: &swarm.SpreadOver{SpreadDescriptor: p.Spread},
				}
				spec.TaskTemplate.Placement.Preferences = append(spec.TaskTemplate.Placement.Preferences, pref)
			}
		}

		// update policy
		spec.UpdateConfig = &swarm.UpdateConfig{
			Parallelism:   info.UpdatePolicy.Parallelism,
			FailureAction: info.UpdatePolicy.FailureAction,
			Order:         info.UpdatePolicy.Order,
		}
		if info.UpdatePolicy.Delay != "" {
			spec.UpdateConfig.Delay, err = time.ParseDuration(info.UpdatePolicy.Delay)
			if err != nil {
				return err
			}
		}

		// rollback policy
		spec.RollbackConfig = &swarm.UpdateConfig{
			Parallelism:   info.RollbackPolicy.Parallelism,
			FailureAction: info.RollbackPolicy.FailureAction,
			Order:         info.RollbackPolicy.Order,
		}
		if info.RollbackPolicy.Delay != "" {
			spec.RollbackConfig.Delay, err = time.ParseDuration(info.RollbackPolicy.Delay)
			if err != nil {
				return err
			}
		}

		// restart policy
		var d time.Duration
		spec.TaskTemplate.RestartPolicy = &swarm.RestartPolicy{
			Condition:   info.RestartPolicy.Condition,
			MaxAttempts: &info.RestartPolicy.MaxAttempts,
		}
		if info.RestartPolicy.Delay != "" {
			d, err = time.ParseDuration(info.RestartPolicy.Delay)
			if err != nil {
				return err
			}
			spec.TaskTemplate.RestartPolicy.Delay = &d
		}
		if info.RestartPolicy.Window != "" {
			d, err = time.ParseDuration(info.RestartPolicy.Window)
			if err != nil {
				return err
			}
			spec.TaskTemplate.RestartPolicy.Window = &d
		}

		// resources
		if info.Resource.Limit.IsSet() || info.Resource.Reserve.IsSet() {
			spec.TaskTemplate.Resources = &swarm.ResourceRequirements{}
			if info.Resource.Limit.IsSet() {
				spec.TaskTemplate.Resources.Limits, err = info.Resource.Limit.ToResources()
				if err != nil {
					return err
				}
			}
			if info.Resource.Limit.IsSet() {
				spec.TaskTemplate.Resources.Reservations, err = info.Resource.Reserve.ToResources()
				if err != nil {
					return err
				}
			}
		}

		// log driver
		if info.LogDriver.Name != "" {
			spec.TaskTemplate.LogDriver = &swarm.Driver{
				Name:    info.LogDriver.Name,
				Options: info.LogDriver.Options.ToMap(),
			}
		}

		// DNS
		spec.TaskTemplate.ContainerSpec.DNSConfig = info.GetDNSConfig()

		options := types.ServiceUpdateOptions{
			RegistryAuthFrom: types.RegistryAuthFromSpec,
			QueryRegistry:    false,
		}
		resp, err := cli.ServiceUpdate(context.Background(), info.Name, version(info.Version), spec, options)
		if err == nil && len(resp.Warnings) > 0 {
			mgr.Logger().Warnf("service %s was updated but got warnings: %v", info.Name, resp.Warnings)
		}
		return err
	})
}

// ServiceScale adjust replicas of a service.
func ServiceScale(name string, version, count uint64) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		service, _, err := cli.ServiceInspectWithRaw(ctx, name, types.ServiceInspectOptions{})
		if err != nil {
			return err
		}

		spec := service.Spec
		if spec.Mode.Replicated == nil {
			return errors.New("the mode of service isn't replicated")
		}
		spec.Mode.Replicated.Replicas = &count

		options := types.ServiceUpdateOptions{
			RegistryAuthFrom: types.RegistryAuthFromSpec,
			QueryRegistry:    false,
		}
		ver := service.Version
		if version > 0 {
			ver = swarm.Version{Index: version}
		}
		resp, err := cli.ServiceUpdate(context.Background(), name, ver, spec, options)
		if err == nil && len(resp.Warnings) > 0 {
			mgr.Logger().Warnf("service %s was scaled but got warnings: %v", name, resp.Warnings)
		}
		return err
	})
}

// ServiceCreate create a service.
func ServiceCreate(info *model.ServiceInfo) error { // nolint: gocyclo
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		service := swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name:   info.Name,
				Labels: info.ServiceLabels.ToMap(),
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: &swarm.ContainerSpec{
					Image:  info.Image,
					Dir:    info.Dir,
					User:   info.User,
					Labels: info.ContainerLabels.ToMap(),
				},
			},
			EndpointSpec: &swarm.EndpointSpec{Mode: swarm.ResolutionModeVIP},
		}

		if info.Command != "" {
			service.TaskTemplate.ContainerSpec.Command = strings.Split(info.Command, " ")
		}
		if info.Args != "" {
			service.TaskTemplate.ContainerSpec.Args = strings.Split(info.Args, " ")
		}

		if info.Mode == "replicated" {
			service.Mode.Replicated = &swarm.ReplicatedService{Replicas: &info.Replicas}
		} else if info.Mode == "global" {
			service.Mode.Global = &swarm.GlobalService{}
		}

		if envs := info.Environments.ToMap(); len(envs) > 0 {
			for n, v := range envs {
				service.TaskTemplate.ContainerSpec.Env = append(service.TaskTemplate.ContainerSpec.Env, n+"="+v)
			}
		}

		for _, n := range info.Networks {
			service.TaskTemplate.Networks = append(service.TaskTemplate.Networks, swarm.NetworkAttachmentConfig{Target: n})
		}

		if info.Endpoint.Mode != "" && len(info.Endpoint.Ports) > 0 {
			service.EndpointSpec = &swarm.EndpointSpec{Mode: swarm.ResolutionMode(info.Endpoint.Mode)}
			for _, p := range info.Endpoint.Ports {
				port := swarm.PortConfig{
					Protocol:      p.Protocol,
					TargetPort:    p.TargetPort,
					PublishedPort: p.PublishedPort,
					PublishMode:   p.PublishMode,
				}
				service.EndpointSpec.Ports = append(service.EndpointSpec.Ports, port)
			}
		}

		for _, s := range info.Secrets {
			service.TaskTemplate.ContainerSpec.Secrets = append(service.TaskTemplate.ContainerSpec.Secrets, s.ToSecret())
		}

		for _, c := range info.Configs {
			service.TaskTemplate.ContainerSpec.Configs = append(service.TaskTemplate.ContainerSpec.Configs, c.ToConfig())
		}

		for _, m := range info.Mounts {
			if m.Target != "" {
				mnt := mount.Mount{
					Type:     mount.Type(m.Type),
					Source:   m.Source,
					Target:   m.Target,
					ReadOnly: m.ReadOnly,
				}
				if m.Propagation != "" {
					mnt.BindOptions = &mount.BindOptions{Propagation: m.Propagation}
				}
				service.TaskTemplate.ContainerSpec.Mounts = append(service.TaskTemplate.ContainerSpec.Mounts, mnt)
			}
		}

		service.TaskTemplate.Placement = &swarm.Placement{}
		for _, c := range info.Placement.Constraints {
			if cons := c.ToConstraint(); cons != "" {
				service.TaskTemplate.Placement.Constraints = append(service.TaskTemplate.Placement.Constraints, cons)
			}
		}
		for _, p := range info.Placement.Preferences {
			if p.Spread != "" {
				pref := swarm.PlacementPreference{
					Spread: &swarm.SpreadOver{SpreadDescriptor: p.Spread},
				}
				service.TaskTemplate.Placement.Preferences = append(service.TaskTemplate.Placement.Preferences, pref)
			}
		}

		// update policy
		service.UpdateConfig = &swarm.UpdateConfig{
			Parallelism:   info.UpdatePolicy.Parallelism,
			FailureAction: info.UpdatePolicy.FailureAction,
			Order:         info.UpdatePolicy.Order,
		}
		if info.UpdatePolicy.Delay != "" {
			service.UpdateConfig.Delay, err = time.ParseDuration(info.UpdatePolicy.Delay)
			if err != nil {
				return err
			}
		}

		// rollback policy
		service.RollbackConfig = &swarm.UpdateConfig{
			Parallelism:   info.RollbackPolicy.Parallelism,
			FailureAction: info.RollbackPolicy.FailureAction,
			Order:         info.RollbackPolicy.Order,
		}
		if info.RollbackPolicy.Delay != "" {
			service.RollbackConfig.Delay, err = time.ParseDuration(info.RollbackPolicy.Delay)
			if err != nil {
				return err
			}
		}

		// restart policy
		var d time.Duration
		service.TaskTemplate.RestartPolicy = &swarm.RestartPolicy{
			Condition:   info.RestartPolicy.Condition,
			MaxAttempts: &info.RestartPolicy.MaxAttempts,
		}
		if info.RestartPolicy.Delay != "" {
			d, err = time.ParseDuration(info.RestartPolicy.Delay)
			if err != nil {
				return err
			}
			service.TaskTemplate.RestartPolicy.Delay = &d
		}
		if info.RestartPolicy.Window != "" {
			d, err = time.ParseDuration(info.RestartPolicy.Window)
			if err != nil {
				return err
			}
			service.TaskTemplate.RestartPolicy.Window = &d
		}

		// resources
		if info.Resource.Limit.IsSet() || info.Resource.Reserve.IsSet() {
			service.TaskTemplate.Resources = &swarm.ResourceRequirements{}
			if info.Resource.Limit.IsSet() {
				service.TaskTemplate.Resources.Limits, err = info.Resource.Limit.ToResources()
				if err != nil {
					return err
				}
			}
			if info.Resource.Limit.IsSet() {
				service.TaskTemplate.Resources.Reservations, err = info.Resource.Reserve.ToResources()
				if err != nil {
					return err
				}
			}
		}

		// log driver
		if info.LogDriver.Name != "" {
			service.TaskTemplate.LogDriver = &swarm.Driver{
				Name:    info.LogDriver.Name,
				Options: info.LogDriver.Options.ToMap(),
			}
		}

		// DNS
		service.TaskTemplate.ContainerSpec.DNSConfig = info.GetDNSConfig()

		opts := types.ServiceCreateOptions{EncodedRegistryAuth: info.RegistryAuth}
		var resp types.ServiceCreateResponse
		resp, err = cli.ServiceCreate(ctx, service, opts)
		if err == nil && len(resp.Warnings) > 0 {
			mgr.Logger().Warnf("service %s was created but got warnings: %v", info.Name, resp.Warnings)
		}
		return
	})
}

// ServiceRemove remove a service.
func ServiceRemove(name string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		return cli.ServiceRemove(ctx, name)
	})
}

// ServiceLogs returns the logs generated by a service.
func ServiceLogs(name string, line int, timestamps bool) (stdout, stderr *bytes.Buffer, err error) {
	var (
		ctx context.Context
		cli *client.Client
		rc  io.ReadCloser
	)

	ctx, cli, err = mgr.Client()
	if err != nil {
		return
	}

	opts := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       strconv.Itoa(line),
		Timestamps: timestamps,
		//Since: (time.Hour * 24).String()
	}
	if rc, err = cli.ServiceLogs(ctx, name, opts); err == nil {
		defer rc.Close()

		stdout = &bytes.Buffer{}
		stderr = &bytes.Buffer{}
		_, err = stdcopy.StdCopy(stdout, stderr, rc)
	}
	return
}

// ServiceCommand returns the docker command line to create this service.
func ServiceCommand(name string) (cmd string, err error) { // nolint: gocyclo
	// todo: support auto wrapping
	var (
		ctx context.Context
		cli *client.Client
	)

	ctx, cli, err = mgr.Client()
	if err != nil {
		return
	}

	service, _, err := cli.ServiceInspectWithRaw(ctx, name, types.ServiceInspectOptions{})
	if err != nil {
		return
	}

	si := model.NewServiceInfo(service)
	b := texts.Builder{}
	b.WriteString("docker service create --name ", si.Name)
	if si.Mode == "global" {
		b.WriteString(" --mode global")
	} else if si.Replicas > 1 {
		b.WriteFormat(" --replicas %d", si.Replicas)
	}
	for _, n := range service.Spec.TaskTemplate.Networks {
		var network types.NetworkResource
		network, err = cli.NetworkInspect(ctx, n.Target, types.NetworkInspectOptions{})
		b.WriteString(" --network ", network.Name)
	}
	if len(si.Endpoint.Ports) > 0 {
		if si.Endpoint.Mode == swarm.ResolutionModeDNSRR {
			b.WriteFormat(" --endpoint-mode %s", si.Endpoint.Mode)
		}
		for _, p := range si.Endpoint.Ports {
			b.WriteFormat(" --publish mode=%s,target=%d,published=%d,protocol=%s", p.PublishMode, p.TargetPort, p.PublishedPort, p.Protocol)
		}
	}
	for _, c := range si.Placement.Constraints {
		b.WriteFormat(" --constraint '%s'", c.ToConstraint())
	}
	for _, p := range si.Placement.Preferences {
		b.WriteFormat(" --placement-pref '%s'", p.Spread)
	}
	for _, e := range si.Environments {
		b.WriteFormat(" --env '%s=%s'", e.Name, e.Value)
	}
	for _, l := range si.ServiceLabels {
		b.WriteFormat(" --label '%s=%s'", l.Name, l.Value)
	}
	for _, l := range si.ContainerLabels {
		b.WriteFormat(" --container-label '%s=%s'", l.Name, l.Value)
	}
	for _, mnt := range si.Mounts {
		// todo: add mnt.Propagation
		b.WriteFormat(" --mount type=%s,source=%s,destination=%s", mnt.Type, mnt.Source, mnt.Target)
		if mnt.ReadOnly {
			b.WriteString(",ro=1")
		}
	}
	// todo: uid/gid
	for _, c := range si.Configs {
		b.WriteFormat(" --config source=%s,target=%s,mode=0%d", c.Name, c.FileName, c.Mode)
	}
	for _, c := range si.Secrets {
		b.WriteFormat(" --secret source=%s,target=%s,mode=0%d", c.Name, c.FileName, c.Mode)
	}
	if si.Resource.Limit.IsSet() {
		if si.Resource.Limit.CPU > 0 {
			b.WriteFormat(" --limit-cpu %.2f", si.Resource.Limit.CPU)
		}
		if si.Resource.Limit.Memory != "" {
			b.WriteString(" --limit-memory ", si.Resource.Limit.Memory)
		}
	}
	if si.Resource.Reserve.IsSet() {
		if si.Resource.Reserve.CPU > 0 {
			b.WriteFormat(" --reserve-cpu %.2f", si.Resource.Reserve.CPU)
		}
		if si.Resource.Reserve.Memory != "" {
			b.WriteString(" --reserve-memory ", si.Resource.Reserve.Memory)
		}
	}
	if si.LogDriver.Name != "" {
		b.WriteString(" --log-driver ", si.LogDriver.Name)
		for _, opt := range si.LogDriver.Options {
			b.WriteFormat(" --log-opt '%s=%s'", opt.Name, opt.Value)
		}
	}
	// UpdatePolicy
	if si.UpdatePolicy.Parallelism > 0 {
		b.WriteFormat(" --update-parallelism %d", si.UpdatePolicy.Parallelism)
	}
	if si.UpdatePolicy.Delay != "" {
		b.WriteString(" --update-delay ", si.UpdatePolicy.Delay)
	}
	if si.UpdatePolicy.FailureAction != "" {
		b.WriteString(" --update-failure-action ", si.UpdatePolicy.FailureAction)
	}
	if si.UpdatePolicy.Order != "" {
		b.WriteString(" --update-order ", si.UpdatePolicy.Order)
	}
	// RollbackPolicy
	if si.RollbackPolicy.Parallelism > 0 {
		b.WriteFormat(" --rollback-parallelism %d", si.RollbackPolicy.Parallelism)
	}
	if si.RollbackPolicy.Delay != "" {
		b.WriteString(" --rollback-delay ", si.RollbackPolicy.Delay)
	}
	if si.RollbackPolicy.FailureAction != "" {
		b.WriteString(" --rollback-failure-action ", si.RollbackPolicy.FailureAction)
	}
	if si.RollbackPolicy.Order != "" {
		b.WriteString(" --rollback-order ", si.RollbackPolicy.Order)
	}
	// RestartPolicy
	if si.RestartPolicy.Condition != "" {
		b.WriteFormat(" --restart-condition %s", si.RestartPolicy.Condition)
	}
	if si.RestartPolicy.MaxAttempts > 0 {
		b.WriteFormat(" --restart-max-attempts %d", si.RestartPolicy.MaxAttempts)
	}
	if si.RestartPolicy.Delay != "" {
		b.WriteString(" --restart-delay ", si.RestartPolicy.Delay)
	}
	if si.RestartPolicy.Window != "" {
		b.WriteString(" --restart-window ", si.RestartPolicy.Window)
	}
	//b.WriteString(" --with-registry-auth")
	if si.Dir != "" {
		b.WriteString(" --workdir ", si.Dir)
	}
	if si.User != "" {
		b.WriteString(" --user ", si.User)
	}
	b.WriteString(" ", si.Image)
	if si.Command != "" {
		b.WriteString(" ", si.Command)
	}
	if si.Args != "" {
		b.WriteString(" ", si.Args)
	}
	cmd = b.String()
	return
}

// ServiceRollback rollbacks a service.
func ServiceRollback(name string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		service, _, err := cli.ServiceInspectWithRaw(ctx, name, types.ServiceInspectOptions{})
		if err != nil {
			return err
		}

		if service.PreviousSpec == nil {
			return errors.New("can't rollback service, previous spec is not exists")
		}

		spec := *service.PreviousSpec
		options := types.ServiceUpdateOptions{
			RegistryAuthFrom: types.RegistryAuthFromPreviousSpec,
			QueryRegistry:    false,
			Rollback:         "previous",
		}
		resp, err := cli.ServiceUpdate(context.Background(), name, service.Version, spec, options)
		if err == nil && len(resp.Warnings) > 0 {
			mgr.Logger().Warnf("service %s was rollbacked but got warnings: %v", name, resp.Warnings)
		}
		return err
	})
}

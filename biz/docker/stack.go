package docker

import (
	"context"
	"fmt"
	"sort"

	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/biz/docker/compose"
	"github.com/cuigh/swirl/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

const stackLabel = "com.docker.stack.namespace"

// StackList return all stacks.
func StackList() (stacks []*model.StackListInfo, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		var services []swarm.Service
		opts := types.ServiceListOptions{
			Filters: filters.NewArgs(),
		}
		opts.Filters.Add("label", stackLabel)
		services, err = cli.ServiceList(ctx, opts)
		if err != nil {
			return
		}

		m := make(map[string]*model.StackListInfo)
		for _, service := range services {
			labels := service.Spec.Labels
			name, ok := labels[stackLabel]
			if !ok {
				err = fmt.Errorf("cannot get label %s for service %s(%s)", stackLabel, service.Spec.Name, service.ID)
				return
			}

			if stack, ok := m[name]; ok {
				stack.Services = append(stack.Services, service.Spec.Name)
			} else {
				m[name] = &model.StackListInfo{
					Name:     name,
					Services: []string{service.Spec.Name},
				}
			}
		}

		for _, stack := range m {
			stacks = append(stacks, stack)
		}
		sort.Slice(stacks, func(i, j int) bool {
			return stacks[i].Name < stacks[j].Name
		})
		return
	})
	return
}

// StackCount return number of stacks.
func StackCount() (count int, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		var services []swarm.Service
		opts := types.ServiceListOptions{
			Filters: filters.NewArgs(),
		}
		opts.Filters.Add("label", stackLabel)
		services, err = cli.ServiceList(ctx, opts)
		if err != nil {
			return
		}

		m := make(map[string]struct{})
		for _, service := range services {
			labels := service.Spec.Labels
			if name, ok := labels[stackLabel]; ok {
				m[name] = struct{}{}
			} else {
				mgr.Logger().Warnf("cannot get label %s for service %s(%s)", stackLabel, service.Spec.Name, service.ID)
			}
		}
		count = len(m)
		return
	})
	return
}

// StackRemove remove a stack.
func StackRemove(name string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		var (
			services []swarm.Service
			networks []types.NetworkResource
			secrets  []swarm.Secret
			configs  []swarm.Config
			errs     []error
		)

		args := filters.NewArgs()
		args.Add("label", stackLabel+"="+name)

		services, err = cli.ServiceList(ctx, types.ServiceListOptions{Filters: args})
		if err != nil {
			return
		}

		networks, err = cli.NetworkList(ctx, types.NetworkListOptions{Filters: args})
		if err != nil {
			return
		}

		// API version >= 1.25
		secrets, err = cli.SecretList(ctx, types.SecretListOptions{Filters: args})
		if err != nil {
			return
		}

		// API version >= 1.30
		configs, err = cli.ConfigList(ctx, types.ConfigListOptions{Filters: args})
		if err != nil {
			return
		}

		if len(services)+len(networks)+len(secrets)+len(configs) == 0 {
			return fmt.Errorf("nothing found in stack: %s", name)
		}

		// Remove services
		for _, service := range services {
			if err = cli.ServiceRemove(ctx, service.ID); err != nil {
				e := errors.Format("Failed to remove service %s: %s", service.Spec.Name, err)
				errs = append(errs, e)
				mgr.Logger().Warn(e)
			}
		}

		// Remove secrets
		for _, secret := range secrets {
			if err = cli.SecretRemove(ctx, secret.ID); err != nil {
				e := errors.Format("Failed to remove secret %s: %s", secret.Spec.Name, err)
				errs = append(errs, e)
				mgr.Logger().Warn(e)
			}
		}

		// Remove configs
		for _, config := range configs {
			if err = cli.ConfigRemove(ctx, config.ID); err != nil {
				e := errors.Format("Failed to remove config %s: %s", config.Spec.Name, err)
				errs = append(errs, e)
				mgr.Logger().Warn(e)
			}
		}

		// Remove networks
		for _, network := range networks {
			if err = cli.NetworkRemove(ctx, network.ID); err != nil {
				e := errors.Format("Failed to remove network %s: %s", network.Name, err)
				errs = append(errs, e)
				mgr.Logger().Warn(e)
			}
		}

		if len(errs) > 0 {
			return errors.List(errs...)
		}
		return nil
	})
}

// StackDeploy deploy a stack.
func StackDeploy(name, content string, authes map[string]string) error {
	ctx, cli, err := mgr.Client()
	if err != nil {
		return err
	}

	cfg, err := compose.Parse(name, content)
	if err != nil {
		return err
	}

	namespace := compose.NewNamespace(name)

	serviceNetworks := compose.GetServicesDeclaredNetworks(cfg.Services)
	networks, externalNetworks := compose.Networks(namespace, cfg.Networks, serviceNetworks)
	if err = validateExternalNetworks(ctx, cli, externalNetworks); err != nil {
		return err
	}
	if err = createNetworks(ctx, cli, namespace, networks); err != nil {
		return err
	}

	secrets, err := compose.Secrets(namespace, cfg.Secrets)
	if err != nil {
		return err
	}
	if err = createSecrets(ctx, cli, secrets); err != nil {
		return err
	}

	configs, err := compose.Configs(namespace, cfg.Configs)
	if err != nil {
		return err
	}
	if err = createConfigs(ctx, cli, configs); err != nil {
		return err
	}

	services, err := compose.Services(namespace, cfg, cli)
	if err != nil {
		return err
	}
	return deployServices(ctx, cli, services, namespace, authes)
}

func validateExternalNetworks(ctx context.Context, cli *client.Client, externalNetworks []string) error {
	for _, networkName := range externalNetworks {
		if !container.NetworkMode(networkName).IsUserDefined() {
			// Networks that are not user defined always exist on all nodes as
			// local-scoped networks, so there's no need to inspect them.
			continue
		}
		network, err := cli.NetworkInspect(ctx, networkName, types.NetworkInspectOptions{})
		switch {
		case client.IsErrNotFound(err):
			return errors.Format("network %q is declared as external, but could not be found. You need to create a swarm-scoped network before the stack is deployed", networkName)
		case err != nil:
			return err
		case network.Scope != "swarm":
			return errors.Format("network %q is declared as external, but it is not in the right scope: %q instead of \"swarm\"", networkName, network.Scope)
		}
	}
	return nil
}

func createNetworks(ctx context.Context, cli *client.Client, namespace compose.Namespace, networks map[string]types.NetworkCreate) error {
	opts := types.NetworkListOptions{
		Filters: filters.NewArgs(),
	}
	opts.Filters.Add("label", stackLabel+"="+namespace.Name())
	existingNetworks, err := cli.NetworkList(ctx, opts)
	if err != nil {
		return err
	}

	existingNetworkMap := make(map[string]types.NetworkResource)
	for _, network := range existingNetworks {
		existingNetworkMap[network.Name] = network
	}

	for internalName, createOpts := range networks {
		name := namespace.Scope(internalName)
		if _, exists := existingNetworkMap[name]; exists {
			continue
		}

		if createOpts.Driver == "" {
			createOpts.Driver = "overlay"
		}

		mgr.Logger().Infof("Creating network %s", name)
		if _, err = cli.NetworkCreate(ctx, name, createOpts); err != nil {
			return errors.Wrap(err, "failed to create network "+internalName)
		}
	}
	return nil
}

func createSecrets(ctx context.Context, cli *client.Client, secrets []swarm.SecretSpec) error {
	for _, secretSpec := range secrets {
		secret, _, err := cli.SecretInspectWithRaw(ctx, secretSpec.Name)
		switch {
		case err == nil:
			// secret already exists, then we update that
			if err = cli.SecretUpdate(ctx, secret.ID, secret.Meta.Version, secretSpec); err != nil {
				return errors.Wrap(err, "failed to update secret "+secretSpec.Name)
			}
		case client.IsErrNotFound(err):
			// secret does not exist, then we create a new one.
			if _, err = cli.SecretCreate(ctx, secretSpec); err != nil {
				return errors.Wrap(err, "failed to create secret "+secretSpec.Name)
			}
		default:
			return err
		}
	}
	return nil
}

func createConfigs(ctx context.Context, cli *client.Client, configs []swarm.ConfigSpec) error {
	for _, configSpec := range configs {
		config, _, err := cli.ConfigInspectWithRaw(ctx, configSpec.Name)
		switch {
		case err == nil:
			// config already exists, then we update that
			if err = cli.ConfigUpdate(ctx, config.ID, config.Meta.Version, configSpec); err != nil {
				errors.Wrap(err, "failed to update config "+configSpec.Name)
			}
		case client.IsErrNotFound(err):
			// config does not exist, then we create a new one.
			if _, err = cli.ConfigCreate(ctx, configSpec); err != nil {
				errors.Wrap(err, "failed to create config "+configSpec.Name)
			}
		default:
			return err
		}
	}
	return nil
}

func getServices(
	ctx context.Context,
	cli *client.Client,
	namespace string,
) ([]swarm.Service, error) {
	opts := types.ServiceListOptions{
		Filters: filters.NewArgs(),
	}
	opts.Filters.Add("label", stackLabel+"="+namespace)
	return cli.ServiceList(ctx, opts)
}

func deployServices(
	ctx context.Context,
	cli *client.Client,
	services map[string]swarm.ServiceSpec,
	namespace compose.Namespace,
	authes map[string]string,
	//sendAuth bool,
	//resolveImage string,
) error {
	existingServices, err := getServices(ctx, cli, namespace.Name())
	if err != nil {
		return err
	}

	existingServiceMap := make(map[string]swarm.Service)
	for _, service := range existingServices {
		existingServiceMap[service.Spec.Name] = service
	}

	for internalName, serviceSpec := range services {
		name := namespace.Scope(internalName)

		// TODO: Add auth
		encodedAuth := authes[serviceSpec.TaskTemplate.ContainerSpec.Image]
		//image := serviceSpec.TaskTemplate.ContainerSpec.Image
		//if sendAuth {
		//	// Retrieve encoded auth token from the image reference
		//	encodedAuth, err = command.RetrieveAuthTokenFromImage(ctx, dockerCli, image)
		//	if err != nil {
		//		return err
		//	}
		//}

		if service, exists := existingServiceMap[name]; exists {
			mgr.Logger().Infof("Updating service %s (id: %s)", name, service.ID)

			updateOpts := types.ServiceUpdateOptions{
				RegistryAuthFrom:    types.RegistryAuthFromSpec,
				EncodedRegistryAuth: encodedAuth,
			}

			//if resolveImage == resolveImageAlways || (resolveImage == resolveImageChanged && image != service.Spec.Labels[compose.LabelImage]) {
			//	updateOpts.QueryRegistry = true
			//}

			response, err := cli.ServiceUpdate(
				ctx,
				service.ID,
				service.Version,
				serviceSpec,
				updateOpts,
			)
			if err != nil {
				return errors.Wrap(err, "failed to update service "+name)
			}

			for _, warning := range response.Warnings {
				mgr.Logger().Warn(warning)
			}
		} else {
			mgr.Logger().Infof("Creating service %s", name)

			createOpts := types.ServiceCreateOptions{EncodedRegistryAuth: encodedAuth}

			// query registry if flag disabling it was not set
			//if resolveImage == resolveImageAlways || resolveImage == resolveImageChanged {
			//	createOpts.QueryRegistry = true
			//}

			if _, err = cli.ServiceCreate(ctx, serviceSpec, createOpts); err != nil {
				return errors.Wrap(err, "failed to create service "+name)
			}
		}
	}
	return nil
}

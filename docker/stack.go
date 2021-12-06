package docker

import (
	"context"

	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/docker/compose"
	composetypes "github.com/cuigh/swirl/docker/compose/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

const stackLabel = "com.docker.stack.namespace"

// StackList return all stacks.
func (d *Docker) StackList(ctx context.Context) (stacks map[string][]string, err error) {
	err = d.call(func(c *client.Client) (err error) {
		var services []swarm.Service
		opts := types.ServiceListOptions{
			Filters: filters.NewArgs(),
		}
		opts.Filters.Add("label", stackLabel)
		services, err = c.ServiceList(ctx, opts)
		if err != nil {
			return
		}

		stacks = make(map[string][]string)
		for _, service := range services {
			name := service.Spec.Labels[stackLabel]
			stacks[name] = append(stacks[name], service.Spec.Name)
		}
		return
	})
	return
}

// StackRemove remove a stack.
func (d *Docker) StackRemove(ctx context.Context, name string) error {
	return d.call(func(c *client.Client) (err error) {
		var (
			services []swarm.Service
			networks []types.NetworkResource
			secrets  []swarm.Secret
			configs  []swarm.Config
			errs     []error
		)

		args := filters.NewArgs()
		args.Add("label", stackLabel+"="+name)

		services, err = c.ServiceList(ctx, types.ServiceListOptions{Filters: args})
		if err != nil {
			return
		}

		networks, err = c.NetworkList(ctx, types.NetworkListOptions{Filters: args})
		if err != nil {
			return
		}

		// API version >= 1.25
		secrets, err = c.SecretList(ctx, types.SecretListOptions{Filters: args})
		if err != nil {
			return
		}

		// API version >= 1.30
		configs, err = c.ConfigList(ctx, types.ConfigListOptions{Filters: args})
		if err != nil {
			return
		}

		if len(services)+len(networks)+len(secrets)+len(configs) == 0 {
			//return fmt.Errorf("nothing found in stack: %s", name)
			return nil
		}

		// Remove services
		for _, service := range services {
			if err = c.ServiceRemove(ctx, service.ID); err != nil {
				e := errors.Format("Failed to remove service %s: %s", service.Spec.Name, err)
				errs = append(errs, e)
				d.logger.Warn(e)
			}
		}

		// Remove secrets
		for _, secret := range secrets {
			if err = c.SecretRemove(ctx, secret.ID); err != nil {
				e := errors.Format("Failed to remove secret %s: %s", secret.Spec.Name, err)
				errs = append(errs, e)
				d.logger.Warn(e)
			}
		}

		// Remove configs
		for _, config := range configs {
			if err = c.ConfigRemove(ctx, config.ID); err != nil {
				e := errors.Format("Failed to remove config %s: %s", config.Spec.Name, err)
				errs = append(errs, e)
				d.logger.Warn(e)
			}
		}

		// Remove networks
		for _, network := range networks {
			if err = c.NetworkRemove(ctx, network.ID); err != nil {
				e := errors.Format("Failed to remove network %s: %s", network.Name, err)
				errs = append(errs, e)
				d.logger.Warn(e)
			}
		}

		if len(errs) > 0 {
			return errors.List(errs...)
		}
		return nil
	})
}

// StackDeploy deploy a stack.
func (d *Docker) StackDeploy(ctx context.Context, cfg *composetypes.Config, authes map[string]string) error {
	c, err := d.client()
	if err != nil {
		return err
	}

	namespace := compose.NewNamespace(cfg.Filename)

	serviceNetworks := compose.GetServicesDeclaredNetworks(cfg.Services)
	networks, externalNetworks := compose.Networks(namespace, cfg.Networks, serviceNetworks)
	if err = validateExternalNetworks(ctx, c, externalNetworks); err != nil {
		return err
	}
	if err = d.createNetworks(ctx, c, namespace, networks); err != nil {
		return err
	}

	secrets, err := compose.Secrets(namespace, cfg.Secrets)
	if err != nil {
		return err
	}
	if err = createSecrets(ctx, c, secrets); err != nil {
		return err
	}

	configs, err := compose.Configs(namespace, cfg.Configs)
	if err != nil {
		return err
	}
	if err = createConfigs(ctx, c, configs); err != nil {
		return err
	}

	services, err := compose.Services(namespace, cfg, c)
	if err != nil {
		return err
	}
	return d.deployServices(ctx, c, services, namespace, authes)
}

// StackCount return number of stacks.
func (d *Docker) StackCount(ctx context.Context) (count int, err error) {
	err = d.call(func(c *client.Client) (err error) {
		var services []swarm.Service
		opts := types.ServiceListOptions{Filters: filters.NewArgs()}
		opts.Filters.Add("label", stackLabel)
		services, err = c.ServiceList(ctx, opts)
		if err != nil {
			return
		}

		m := make(map[string]struct{})
		for _, service := range services {
			labels := service.Spec.Labels
			m[labels[stackLabel]] = struct{}{}
		}
		count = len(m)
		return
	})
	return
}

func validateExternalNetworks(ctx context.Context, c *client.Client, externalNetworks []string) error {
	for _, networkName := range externalNetworks {
		if !container.NetworkMode(networkName).IsUserDefined() {
			// Networks that are not user defined always exist on all nodes as
			// local-scoped networks, so there's no need to inspect them.
			continue
		}
		network, err := c.NetworkInspect(ctx, networkName, types.NetworkInspectOptions{})
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

func (d *Docker) createNetworks(ctx context.Context, c *client.Client, namespace compose.Namespace, networks map[string]types.NetworkCreate) error {
	opts := types.NetworkListOptions{
		Filters: filters.NewArgs(),
	}
	opts.Filters.Add("label", stackLabel+"="+namespace.Name())
	existingNetworks, err := c.NetworkList(ctx, opts)
	if err != nil {
		return err
	}

	existingNetworkMap := make(map[string]types.NetworkResource)
	for _, network := range existingNetworks {
		existingNetworkMap[network.Name] = network
	}

	for name, createOpts := range networks {
		if _, exists := existingNetworkMap[name]; exists {
			continue
		}

		if createOpts.Driver == "" {
			createOpts.Driver = "overlay"
		}

		d.logger.Infof("Creating network %s", name)
		if _, err = c.NetworkCreate(ctx, name, createOpts); err != nil {
			return errors.Wrap(err, "failed to create network "+name)
		}
	}
	return nil
}

func createSecrets(ctx context.Context, c *client.Client, secrets []swarm.SecretSpec) error {
	for _, secretSpec := range secrets {
		secret, _, err := c.SecretInspectWithRaw(ctx, secretSpec.Name)
		switch {
		case err == nil:
			// secret already exists, then we update that
			if err = c.SecretUpdate(ctx, secret.ID, secret.Meta.Version, secretSpec); err != nil {
				return errors.Wrap(err, "failed to update secret "+secretSpec.Name)
			}
		case client.IsErrNotFound(err):
			// secret does not exist, then we create a new one.
			if _, err = c.SecretCreate(ctx, secretSpec); err != nil {
				return errors.Wrap(err, "failed to create secret "+secretSpec.Name)
			}
		default:
			return err
		}
	}
	return nil
}

func createConfigs(ctx context.Context, c *client.Client, configs []swarm.ConfigSpec) error {
	for _, configSpec := range configs {
		config, _, err := c.ConfigInspectWithRaw(ctx, configSpec.Name)
		switch {
		case err == nil:
			// config already exists, then we update that
			if err = c.ConfigUpdate(ctx, config.ID, config.Meta.Version, configSpec); err != nil {
				return errors.Wrap(err, "failed to update config "+configSpec.Name)
			}
		case client.IsErrNotFound(err):
			// config does not exist, then we create a new one.
			if _, err = c.ConfigCreate(ctx, configSpec); err != nil {
				return errors.Wrap(err, "failed to create config "+configSpec.Name)
			}
		default:
			return err
		}
	}
	return nil
}

func getServices(ctx context.Context, c *client.Client, namespace string) ([]swarm.Service, error) {
	opts := types.ServiceListOptions{
		Filters: filters.NewArgs(),
	}
	opts.Filters.Add("label", stackLabel+"="+namespace)
	return c.ServiceList(ctx, opts)
}

func (d *Docker) deployServices(
	ctx context.Context,
	c *client.Client,
	services map[string]swarm.ServiceSpec,
	namespace compose.Namespace,
	authes map[string]string,
	//sendAuth bool,
	//resolveImage string,
) error {
	existingServices, err := getServices(ctx, c, namespace.Name())
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
			d.logger.Infof("Updating service %s (id: %s)", name, service.ID)

			updateOpts := types.ServiceUpdateOptions{
				RegistryAuthFrom:    types.RegistryAuthFromSpec,
				EncodedRegistryAuth: encodedAuth,
			}

			//if resolveImage == resolveImageAlways || (resolveImage == resolveImageChanged && image != service.Spec.Labels[compose.LabelImage]) {
			//	updateOpts.QueryRegistry = true
			//}

			response, err := c.ServiceUpdate(
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
				d.logger.Warn(warning)
			}
		} else {
			d.logger.Infof("Creating service %s", name)

			createOpts := types.ServiceCreateOptions{EncodedRegistryAuth: encodedAuth}

			// query registry if flag disabling it was not set
			//if resolveImage == resolveImageAlways || resolveImage == resolveImageChanged {
			//	createOpts.QueryRegistry = true
			//}

			if _, err = c.ServiceCreate(ctx, serviceSpec, createOpts); err != nil {
				return errors.Wrap(err, "failed to create service "+name)
			}
		}
	}
	return nil
}

package docker

import (
	"context"
	"sort"

	"github.com/cuigh/swirl/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// NetworkList return all networks.
func NetworkList() (networks []types.NetworkResource, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		networks, err = cli.NetworkList(ctx, types.NetworkListOptions{})
		if err == nil {
			sort.Slice(networks, func(i, j int) bool {
				return networks[i].Name < networks[j].Name
			})
		}
		return
	})
	return
}

// NetworkCount return number of networks.
func NetworkCount() (count int, err error) {
	err = mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		var networks []types.NetworkResource
		networks, err = cli.NetworkList(ctx, types.NetworkListOptions{})
		if err == nil {
			count = len(networks)
		}
		return
	})
	return
}

// NetworkCreate create a network.
func NetworkCreate(info *model.NetworkCreateInfo) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		var (
			resp    types.NetworkCreateResponse
			options = types.NetworkCreate{
				Internal:   info.Internal,
				Attachable: info.Attachable,
				IPAM:       &network.IPAM{},
				EnableIPv6: info.IPV6.Enabled,
				Options:    info.Options.ToMap(),
				Labels:     info.Labels.ToMap(),
			}
		)

		if info.Driver == "other" {
			options.Driver = info.CustomDriver
		} else {
			options.Driver = info.Driver
		}

		if info.IPV4.Subnet != "" || info.IPV4.Gateway != "" || info.IPV4.IPRange != "" {
			cfg := network.IPAMConfig{
				Subnet:  info.IPV4.Subnet,
				Gateway: info.IPV4.Gateway,
				IPRange: info.IPV4.IPRange,
			}
			options.IPAM.Config = append(options.IPAM.Config, cfg)
		}

		if info.IPV6.Enabled && (info.IPV6.Subnet != "" || info.IPV6.Gateway != "") {
			cfg := network.IPAMConfig{
				Subnet:  info.IPV6.Subnet,
				Gateway: info.IPV6.Gateway,
			}
			options.IPAM.Config = append(options.IPAM.Config, cfg)
		}

		resp, err = cli.NetworkCreate(ctx, info.Name, options)
		if err == nil && resp.Warning != "" {
			mgr.Logger().Warnf("network %s was created but got warning: %s", info.Name, resp.Warning)
		}
		return
	})
}

// NetworkRemove remove a network.
func NetworkRemove(name string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		return cli.NetworkRemove(ctx, name)
	})
}

// NetworkDisconnect Disconnect a container from a network.
func NetworkDisconnect(name, container string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		return cli.NetworkDisconnect(ctx, name, container, false)
	})
}

// NetworkInspect return network information.
func NetworkInspect(name string) (network types.NetworkResource, err error) {
	var (
		ctx context.Context
		cli *client.Client
	)
	if ctx, cli, err = mgr.Client(); err == nil {
		network, err = cli.NetworkInspect(ctx, name, types.NetworkInspectOptions{})
	}
	return
}

// NetworkInspectRaw return network raw information.
func NetworkInspectRaw(name string) (raw []byte, err error) {
	var (
		ctx context.Context
		cli *client.Client
	)
	if ctx, cli, err = mgr.Client(); err == nil {
		_, raw, err = cli.NetworkInspectWithRaw(ctx, name, types.NetworkInspectOptions{})
	}
	return
}

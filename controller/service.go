package controller

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/cuigh/auxo/data/set"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/model"
)

type ServiceController struct {
	List   web.HandlerFunc `path:"/" name:"service.list" authorize:"!" desc:"service list page"`
	Detail web.HandlerFunc `path:"/:name/detail" name:"service.detail" authorize:"!" desc:"service detail page"`
	Raw    web.HandlerFunc `path:"/:name/raw" name:"service.raw" authorize:"!" desc:"service raw page"`
	Logs   web.HandlerFunc `path:"/:name/logs" name:"service.logs" authorize:"!" desc:"service logs page"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"service.delete" authorize:"!" desc:"delete service"`
	Scale  web.HandlerFunc `path:"/scale" method:"post" name:"service.scale" authorize:"!" desc:"scale service"`
	New    web.HandlerFunc `path:"/new" name:"service.new" authorize:"!" desc:"new service page"`
	Create web.HandlerFunc `path:"/new" method:"post" name:"service.create" authorize:"!" desc:"create service"`
	Edit   web.HandlerFunc `path:"/:name/edit" name:"service.edit" authorize:"!" desc:"service edit page"`
	Update web.HandlerFunc `path:"/:name/update" method:"post" name:"service.update" authorize:"!" desc:"update service"`
}

func Service() (c *ServiceController) {
	c = &ServiceController{}

	c.List = func(ctx web.Context) error {
		name := ctx.Q("name")
		page := cast.ToIntD(ctx.Q("page"), 1)
		services, totalCount, err := docker.ServiceList(name, page, model.PageSize)
		if err != nil {
			return err
		}

		m := newPagerModel(ctx, totalCount, model.PageSize, page).
			Add("Name", name).
			Add("Services", services)
		return ctx.Render("service/list", m)
	}

	c.Detail = func(ctx web.Context) error {
		name := ctx.P("name")
		service, _, err := docker.ServiceInspect(name)
		if err != nil {
			return err
		}

		info := model.NewServiceDetailInfo(service)
		for _, vip := range service.Endpoint.VirtualIPs {
			n, err := docker.NetworkInspect(vip.NetworkID)
			if err != nil {
				return err
			}
			info.Networks = append(info.Networks, model.Network{ID: vip.NetworkID, Name: n.Name, Address: vip.Addr})
		}

		tasks, err := docker.TaskList(name, "")
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Service", info).Add("Tasks", tasks)
		return ctx.Render("service/detail", m)
	}

	c.Raw = func(ctx web.Context) error {
		name := ctx.P("name")
		_, raw, err := docker.ServiceInspect(name)
		if err != nil {
			return err
		}

		buf := &bytes.Buffer{}
		err = json.Indent(buf, raw, "", "    ")
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Service", name).Add("Raw", string(buf.Bytes()))
		return ctx.Render("service/raw", m)
	}

	c.Logs = func(ctx web.Context) error {
		name := ctx.P("name")
		line := cast.ToIntD(ctx.Q("line"), 500)
		timestamps := cast.ToBoolD(ctx.Q("timestamps"), false)
		stdout, stderr, err := docker.ServiceLogs(name, line, timestamps)
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Service", name).Add("Line", line).Add("Timestamps", timestamps).
			Add("Stdout", stdout.String()).Add("Stderr", stderr.String())
		return ctx.Render("service/logs", m)
	}

	c.Delete = func(ctx web.Context) error {
		names := strings.Split(ctx.F("names"), ",")
		for _, name := range names {
			if err := docker.ServiceRemove(name); err != nil {
				return ajaxResult(ctx, err)
			} else {
				biz.Event.CreateService(model.EventActionDelete, name, ctx.User())
			}
		}
		return ajaxSuccess(ctx, nil)
	}

	c.New = func(ctx web.Context) error {
		networks, err := docker.NetworkList()
		if err != nil {
			return err
		}
		secrets, _, err := docker.SecretList("", 1, 100)
		if err != nil {
			return err
		}
		configs, _, err := docker.ConfigList("", 1, 100)
		if err != nil {
			return err
		}
		registries, err := biz.Registry.List()
		if err != nil {
			return err
		}
		m := newModel(ctx).Add("Networks", networks).Add("Secrets", secrets).Add("Configs", configs).Add("Registries", registries)
		return ctx.Render("service/new", m)
	}

	c.Create = func(ctx web.Context) error {
		info := &model.ServiceInfo{}
		err := ctx.Bind(info)
		if err == nil {
			if info.Registry != "" {
				registry, err := biz.Registry.Get(info.Registry)
				if err != nil {
					return errors.Wrap("Load registry info failed", err)
				} else if registry == nil {
					return errors.New("Can't load registry info")
				}

				info.Image = registry.URL + "/" + info.Image
				info.RegistryAuth = registry.GetEncodedAuth()
			}
			err = docker.ServiceCreate(info)
		}

		if err == nil {
			biz.Event.CreateService(model.EventActionCreate, info.Name, ctx.User())
		}
		return ajaxResult(ctx, err)
	}

	c.Edit = func(ctx web.Context) error {
		name := ctx.P("name")
		service, _, err := docker.ServiceInspect(name)
		if err != nil {
			return err
		}

		networks, err := docker.NetworkList()
		if err != nil {
			return err
		}
		secrets, _, err := docker.SecretList("", 1, 100)
		if err != nil {
			return err
		}
		configs, _, err := docker.ConfigList("", 1, 100)
		if err != nil {
			return err
		}

		checkedNetworks := set.FromSlice(service.Endpoint.VirtualIPs, func(i int) interface{} { return service.Endpoint.VirtualIPs[i].NetworkID })
		checkedSecrets := set.FromSlice(service.Spec.TaskTemplate.ContainerSpec.Secrets, func(i int) interface{} {
			return service.Spec.TaskTemplate.ContainerSpec.Secrets[i].SecretName
		})
		checkedConfigs := set.FromSlice(service.Spec.TaskTemplate.ContainerSpec.Configs, func(i int) interface{} {
			return service.Spec.TaskTemplate.ContainerSpec.Configs[i].ConfigName
		})

		m := newModel(ctx).Add("Service", model.NewServiceInfo(service)).
			Add("Networks", networks).Add("CheckedNetworks", checkedNetworks).
			Add("Secrets", secrets).Add("CheckedSecrets", checkedSecrets).
			Add("Configs", configs).Add("CheckedConfigs", checkedConfigs)
		return ctx.Render("service/edit", m)
	}

	c.Update = func(ctx web.Context) error {
		info := &model.ServiceInfo{}
		err := ctx.Bind(info)
		if err == nil {
			info.Name = ctx.P("name")
			err = docker.ServiceUpdate(info)
		}

		if err == nil {
			biz.Event.CreateService(model.EventActionUpdate, info.Name, ctx.User())
		}
		return ajaxResult(ctx, err)
	}

	c.Scale = func(ctx web.Context) error {
		name := ctx.F("name")
		count, err := strconv.Atoi(ctx.F("count"))
		if err != nil {
			return err
		}

		err = docker.ServiceScale(name, uint64(count))
		if err == nil {
			biz.Event.CreateService(model.EventActionScale, name, ctx.User())
		}
		return ajaxResult(ctx, err)
	}

	return
}

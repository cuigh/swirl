package controller

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/cuigh/auxo/data/set"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

// ServiceController is a controller of docker service
type ServiceController struct {
	List     web.HandlerFunc `path:"/" name:"service.list" authorize:"!" desc:"service list page"`
	Detail   web.HandlerFunc `path:"/:name/detail" name:"service.detail" authorize:"!" desc:"service detail page"`
	Raw      web.HandlerFunc `path:"/:name/raw" name:"service.raw" authorize:"!" desc:"service raw page"`
	Logs     web.HandlerFunc `path:"/:name/logs" name:"service.logs" authorize:"!" desc:"service logs page"`
	Delete   web.HandlerFunc `path:"/delete" method:"post" name:"service.delete" authorize:"!" desc:"delete service"`
	Scale    web.HandlerFunc `path:"/scale" method:"post" name:"service.scale" authorize:"!" desc:"scale service"`
	Rollback web.HandlerFunc `path:"/rollback" method:"post" name:"service.rollback" authorize:"!" desc:"rollback service"`
	New      web.HandlerFunc `path:"/new" name:"service.new" authorize:"!" desc:"new service page"`
	Create   web.HandlerFunc `path:"/new" method:"post" name:"service.create" authorize:"!" desc:"create service"`
	Edit     web.HandlerFunc `path:"/:name/edit" name:"service.edit" authorize:"!" desc:"service edit page"`
	Update   web.HandlerFunc `path:"/:name/edit" method:"post" name:"service.update" authorize:"!" desc:"update service"`
}

// Service creates an instance of ServiceController
func Service() (c *ServiceController) {
	return &ServiceController{
		List:     serviceList,
		Detail:   serviceDetail,
		Raw:      serviceRaw,
		Logs:     serviceLogs,
		Delete:   serviceDelete,
		New:      serviceNew,
		Create:   serviceCreate,
		Edit:     serviceEdit,
		Update:   serviceUpdate,
		Scale:    serviceScale,
		Rollback: serviceRollback,
	}
}

func serviceList(ctx web.Context) error {
	name := ctx.Q("name")
	page := cast.ToInt(ctx.Q("page"), 1)
	services, totalCount, err := docker.ServiceList(name, page, model.PageSize)
	if err != nil {
		return err
	}

	m := newPagerModel(ctx, totalCount, model.PageSize, page).
		Add("Name", name).
		Add("Services", services)
	return ctx.Render("service/list", m)
}

func serviceDetail(ctx web.Context) error {
	name := ctx.P("name")
	service, _, err := docker.ServiceInspect(name)
	if err != nil {
		return err
	}

	info := model.NewServiceDetailInfo(service)
	for _, vip := range service.Endpoint.VirtualIPs {
		n, e := docker.NetworkInspect(vip.NetworkID)
		if e != nil {
			return e
		}
		info.Networks = append(info.Networks, model.Network{ID: vip.NetworkID, Name: n.Name, Address: vip.Addr})
	}

	tasks, _, err := docker.TaskList(&model.TaskListArgs{Service: name})
	if err != nil {
		return err
	}

	cmd, err := docker.ServiceCommand(name)
	if err != nil {
		return err
	}

	m := newModel(ctx).Add("Service", info).Add("Tasks", tasks).Add("Command", cmd)
	return ctx.Render("service/detail", m)
}

func serviceRaw(ctx web.Context) error {
	name := ctx.P("name")
	_, raw, err := docker.ServiceInspect(name)
	if err != nil {
		return err
	}

	j, err := misc.JSONIndent(raw)
	if err != nil {
		return err
	}

	m := newModel(ctx).Add("Service", name).Add("Raw", j)
	return ctx.Render("service/raw", m)
}

func serviceLogs(ctx web.Context) error {
	name := ctx.P("name")
	line := cast.ToInt(ctx.Q("line"), 500)
	timestamps := cast.ToBool(ctx.Q("timestamps"), false)
	stdout, stderr, err := docker.ServiceLogs(name, line, timestamps)
	if err != nil {
		return err
	}

	m := newModel(ctx).Add("Service", name).Add("Line", line).Add("Timestamps", timestamps).
		Add("Stdout", stdout.String()).Add("Stderr", stderr.String())
	return ctx.Render("service/logs", m)
}

func serviceDelete(ctx web.Context) error {
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

func serviceNew(ctx web.Context) error {
	service := &model.ServiceInfo{}
	tid := ctx.Q("template")
	if tid != "" {
		tpl, err := biz.Template.Get(tid)
		if err != nil {
			return err
		}

		if tpl != nil {
			err = json.Unmarshal([]byte(tpl.Content), service)
			if err != nil {
				return err
			}

			if service.Registry != "" {
				var registry *model.Registry
				registry, err = biz.Registry.Get(service.Registry)
				if err != nil {
					return err
				}
				service.RegistryURL = registry.URL
			}
		}
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
	registries, err := biz.Registry.List()
	if err != nil {
		return err
	}
	checkedNetworks := set.FromSlice(service.Networks, func(i int) interface{} { return service.Networks[i] })

	m := newModel(ctx).Add("Service", service).Add("Registries", registries).
		Add("Networks", networks).Add("CheckedNetworks", checkedNetworks).
		Add("Secrets", secrets).Add("Configs", configs)
	return ctx.Render("service/new", m)
}

func serviceCreate(ctx web.Context) error {
	info := &model.ServiceInfo{}
	err := ctx.Bind(info)
	if err != nil {
		return err
	}

	if info.Registry != "" {
		var registry *model.Registry
		registry, err = biz.Registry.Get(info.Registry)
		if err != nil {
			return errors.Wrap("Load registry info failed", err)
		} else if registry == nil {
			return errors.New("Can't load registry info")
		}

		info.Image = registry.URL + "/" + info.Image
		info.RegistryAuth = registry.GetEncodedAuth()
	}

	if err = docker.ServiceCreate(info); err == nil {
		biz.Event.CreateService(model.EventActionCreate, info.Name, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func serviceEdit(ctx web.Context) error {
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

	stack := service.Spec.Labels["com.docker.stack.namespace"]
	checkedNetworks := set.FromSlice(service.Endpoint.VirtualIPs, func(i int) interface{} { return service.Endpoint.VirtualIPs[i].NetworkID })

	m := newModel(ctx).Add("Service", model.NewServiceInfo(service)).Add("Stack", stack).
		Add("Networks", networks).Add("CheckedNetworks", checkedNetworks).
		Add("Secrets", secrets).Add("Configs", configs)
	return ctx.Render("service/edit", m)
}

func serviceUpdate(ctx web.Context) error {
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

func serviceScale(ctx web.Context) error {
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

func serviceRollback(ctx web.Context) error {
	name := ctx.F("name")
	err := docker.ServiceRollback(name)
	if err == nil {
		biz.Event.CreateService(model.EventActionRollback, name, ctx.User())
	}
	return ajaxResult(ctx, err)
}

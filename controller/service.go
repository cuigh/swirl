package controller

import (
	"strconv"
	"strings"
	"time"

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
	List       web.HandlerFunc `path:"/" name:"service.list" authorize:"!" desc:"service list page"`
	Detail     web.HandlerFunc `path:"/:name/detail" name:"service.detail" authorize:"!" perm:"read,service,name"`
	Raw        web.HandlerFunc `path:"/:name/raw" name:"service.raw" authorize:"!" perm:"read,service,name"`
	Logs       web.HandlerFunc `path:"/:name/logs" name:"service.logs" authorize:"!" perm:"read,service,name"`
	Delete     web.HandlerFunc `path:"/:name/delete" method:"post" name:"service.delete" authorize:"!" perm:"write,service,name"`
	Scale      web.HandlerFunc `path:"/:name/scale" method:"post" name:"service.scale" authorize:"!" perm:"write,service,name"`
	Rollback   web.HandlerFunc `path:"/:name/rollback" method:"post" name:"service.rollback" authorize:"!" perm:"write,service,name"`
	Restart    web.HandlerFunc `path:"/:name/restart" method:"post" name:"service.restart" authorize:"!" perm:"write,service,name"`
	New        web.HandlerFunc `path:"/new" name:"service.new" authorize:"!" desc:"new service page"`
	Create     web.HandlerFunc `path:"/new" method:"post" name:"service.create" authorize:"!" desc:"create service"`
	Edit       web.HandlerFunc `path:"/:name/edit" name:"service.edit" authorize:"!" perm:"write,service,name"`
	Update     web.HandlerFunc `path:"/:name/edit" method:"post" name:"service.update" authorize:"!" perm:"write,service,name"`
	PermEdit   web.HandlerFunc `path:"/:name/perm" name:"service.perm.edit" authorize:"!" perm:"write,service,name"`
	PermUpdate web.HandlerFunc `path:"/:name/perm" method:"post" name:"service.perm.update" authorize:"!" perm:"write,service,name"`
	Stats      web.HandlerFunc `path:"/:name/stats" name:"service.stats" authorize:"!" perm:"read,service,name"`
}

// Service creates an instance of ServiceController
func Service() (c *ServiceController) {
	return &ServiceController{
		List:       serviceList,
		Detail:     serviceDetail,
		Raw:        serviceRaw,
		Logs:       serviceLogs,
		Delete:     serviceDelete,
		New:        serviceNew,
		Create:     serviceCreate,
		Edit:       serviceEdit,
		Update:     serviceUpdate,
		Scale:      serviceScale,
		Rollback:   serviceRollback,
		Restart:    serviceRestart,
		PermEdit:   servicePermEdit,
		PermUpdate: permUpdate("service", "name"),
		Stats:      serviceStats,
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
		Set("Name", name).
		Set("Services", services)
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

	m := newModel(ctx).Set("Service", info).Set("Tasks", tasks).Set("Command", cmd)
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

	m := newModel(ctx).Set("Service", name).Set("Raw", j)
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

	m := newModel(ctx).Set("Service", name).Set("Line", line).Set("Timestamps", timestamps).
		Set("Stdout", stdout.String()).Set("Stderr", stderr.String())
	return ctx.Render("service/logs", m)
}

func serviceDelete(ctx web.Context) error {
	names := strings.Split(ctx.F("names"), ",")
	for _, name := range names {
		if err := docker.ServiceRemove(name); err != nil {
			return ajaxResult(ctx, err)
		}
		biz.Event.CreateService(model.EventActionDelete, name, ctx.User())
	}
	return ajaxSuccess(ctx, nil)
}

func serviceNew(ctx web.Context) error {
	info := &model.ServiceInfo{}
	tid := ctx.Q("template")
	if tid != "" {
		err := biz.Template.FillInfo(tid, info)
		if err != nil {
			return err
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

	checkedNetworks := set.NewStringSet(info.Networks...)
	m := newModel(ctx).Set("Service", info).Set("Registries", registries).
		Set("Networks", networks).Set("CheckedNetworks", checkedNetworks).
		Set("Secrets", secrets).Set("Configs", configs)
	return ctx.Render("service/new", m)
}

func serviceCreate(ctx web.Context) error {
	info := &model.ServiceInfo{}
	err := ctx.Bind(info)
	if err != nil {
		return err
	}
	info.Normalize()

	if info.Registry != "" {
		var registry *model.Registry
		registry, err = biz.Registry.Get(info.Registry)
		if err != nil {
			return errors.Wrap(err, "load registry info failed")
		} else if registry == nil {
			return errors.New("can't load registry info")
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

	si := model.NewServiceInfo(service)
	names, err := docker.NetworkNames(si.Networks...)
	if err != nil {
		return err
	}
	checkedNetworks := set.NewStringSet(names...)

	m := newModel(ctx).Set("Service", si).
		Set("Stack", service.Spec.Labels["com.docker.stack.namespace"]).
		Set("Networks", networks).Set("CheckedNetworks", checkedNetworks).
		Set("Secrets", secrets).Set("Configs", configs)
	return ctx.Render("service/edit", m)
}

func serviceUpdate(ctx web.Context) error {
	info := &model.ServiceInfo{}
	err := ctx.Bind(info)
	if err == nil {
		info.Name = ctx.P("name")
		info.Normalize()
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

	err = docker.ServiceScale(name, 0, uint64(count))
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

func serviceRestart(ctx web.Context) error {
	name := ctx.F("name")
	err := docker.ServiceRestart(name)
	if err == nil {
		biz.Event.CreateService(model.EventActionRestart, name, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func servicePermEdit(ctx web.Context) error {
	name := ctx.P("name")
	m := newModel(ctx).Set("Name", name)
	return permEdit(ctx, "service", name, "service/perm", m)
}

func serviceStats(ctx web.Context) error {
	name := ctx.P("name")
	service, _, err := docker.ServiceInspect(name)
	if err != nil {
		return err
	}

	tasks, _, err := docker.TaskList(&model.TaskListArgs{Service: name})
	if err != nil {
		return err
	}

	setting, err := biz.Setting.Get()
	if err != nil {
		return err
	}

	var charts []*model.Chart
	if setting.Metrics.Prometheus != "" {
		var dashboard *model.ChartDashboard
		if dashboard, err = biz.Chart.GetDashboard("service", name); err != nil {
			return err
		}

		if dashboard == nil {
			charts, err = biz.Chart.GetServiceCharts(name)
		} else {
			charts, err = biz.Chart.GetDashboardCharts(dashboard)
		}

		if err != nil {
			return err
		}
	}

	period := cast.ToDuration(ctx.Q("time"), time.Hour)
	refresh := cast.ToBool(ctx.Q("refresh"), true)
	m := newModel(ctx).Set("Service", service).Set("Tasks", tasks).Set("Time", period.String()).
		Set("Refresh", refresh).Set("Charts", charts)
	return ctx.Render("service/stats", m)
}

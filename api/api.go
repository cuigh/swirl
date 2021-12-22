package api

import (
	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/net/web"
)

func ajax(ctx web.Context, err error) error {
	if err != nil {
		return err
	}
	return success(ctx, nil)
}

func success(ctx web.Context, data interface{}) error {
	return ctx.Result(0, "", data)
}

func init() {
	container.Put(NewSystem, container.Name("api.system"))
	container.Put(NewSetting, container.Name("api.setting"))
	container.Put(NewUser, container.Name("api.user"))
	container.Put(NewNode, container.Name("api.node"))
	container.Put(NewRegistry, container.Name("api.registry"))
	container.Put(NewNetwork, container.Name("api.network"))
	container.Put(NewService, container.Name("api.service"))
	container.Put(NewTask, container.Name("api.task"))
	container.Put(NewConfig, container.Name("api.config"))
	container.Put(NewSecret, container.Name("api.secret"))
	container.Put(NewStack, container.Name("api.stack"))
	container.Put(NewImage, container.Name("api.image"))
	container.Put(NewContainer, container.Name("api.container"))
	container.Put(NewVolume, container.Name("api.volume"))
	container.Put(NewUser, container.Name("api.user"))
	container.Put(NewRole, container.Name("api.role"))
	container.Put(NewEvent, container.Name("api.event"))
	container.Put(NewChart, container.Name("api.chart"))
	container.Put(NewDashboard, container.Name("api.dashboard"))
}

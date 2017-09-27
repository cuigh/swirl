package controller

import (
	"strings"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/model"
)

type ConfigController struct {
	List   web.HandlerFunc `path:"/" name:"config.list" authorize:"!" desc:"config list page"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"config.delete" authorize:"!" desc:"delete config"`
	New    web.HandlerFunc `path:"/new" name:"config.new" authorize:"!" desc:"new config page"`
	Create web.HandlerFunc `path:"/new" method:"post" name:"config.create" authorize:"!" desc:"create config"`
	Edit   web.HandlerFunc `path:"/:id/edit" name:"config.edit" authorize:"!" desc:"edit config page"`
	Update web.HandlerFunc `path:"/:id/update" method:"post" name:"config.update" authorize:"!" desc:"update config"`
}

func Config() (c *ConfigController) {
	c = &ConfigController{}

	c.List = func(ctx web.Context) error {
		name := ctx.Q("name")
		page := cast.ToIntD(ctx.Q("page"), 1)
		configs, totalCount, err := docker.ConfigList(name, page, model.PageSize)
		if err != nil {
			return err
		}

		m := newPagerModel(ctx, totalCount, model.PageSize, page).
			Add("Name", name).
			Add("Configs", configs)
		return ctx.Render("config/list", m)
	}

	c.Delete = func(ctx web.Context) error {
		ids := strings.Split(ctx.F("ids"), ",")
		err := docker.ConfigRemove(ids)
		return ajaxResult(ctx, err)
	}

	c.New = func(ctx web.Context) error {
		m := newModel(ctx)
		return ctx.Render("config/new", m)
	}

	c.Create = func(ctx web.Context) error {
		v := &model.ConfigCreateInfo{}
		err := ctx.Bind(v)
		if err == nil {
			err = docker.ConfigCreate(v)
			if err == nil {
				biz.Event.CreateConfig(model.EventActionCreate, v.Name, ctx.User())
			}
		}
		return ajaxResult(ctx, err)
	}

	c.Edit = func(ctx web.Context) error {
		id := ctx.P("id")
		cfg, _, err := docker.ConfigInspect(id)
		if err != nil {
			return err
		}
		m := newModel(ctx).Add("Config", cfg)
		return ctx.Render("config/edit", m)
	}

	c.Update = func(ctx web.Context) error {
		v := &model.ConfigUpdateInfo{}
		err := ctx.Bind(v)
		if err == nil {
			err = docker.ConfigUpdate(v)
			if err == nil {
				biz.Event.CreateConfig(model.EventActionUpdate, v.Name, ctx.User())
			}
		}
		return ajaxResult(ctx, err)
	}

	return
}

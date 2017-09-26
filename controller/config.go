package controller

import (
	"strings"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/model"
)

type ConfigController struct {
	List   web.HandlerFunc `path:"/" name:"config.list" authorize:"!" desc:"config list page"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"config.delete" authorize:"!" desc:"delete config"`
	New    web.HandlerFunc `path:"/new" name:"config.new" authorize:"!" desc:"new config page"`
	Create web.HandlerFunc `path:"/new" method:"post" name:"config.create" authorize:"!" desc:"create config"`
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
		v := struct {
			Name   string `json:"name"`
			Data   string `json:"data"`
			Labels []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"labels"`
		}{}
		err := ctx.Bind(&v)
		if err == nil {
			labels := make(map[string]string)
			for _, l := range v.Labels {
				if l.Name != "" && l.Value != "" {
					labels[l.Name] = l.Value
				}
			}
			err = docker.ConfigCreate(v.Name, []byte(v.Data), labels)
		}
		return ajaxResult(ctx, err)
	}

	return
}

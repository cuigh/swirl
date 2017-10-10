package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

type ContainerController struct {
	List   web.HandlerFunc `path:"/" name:"container.list" authorize:"!" desc:"container list page"`
	Detail web.HandlerFunc `path:"/:id/detail" name:"container.detail" authorize:"!" desc:"container detail page"`
	Raw    web.HandlerFunc `path:"/:id/raw" name:"container.raw" authorize:"!" desc:"container raw page"`
}

func Container() (c *ContainerController) {
	c = &ContainerController{}

	c.List = func(ctx web.Context) error {
		name := ctx.Q("name")
		page := cast.ToIntD(ctx.Q("page"), 1)
		containers, totalCount, err := docker.ContainerList(name, page, model.PageSize)
		if err != nil {
			return err
		}

		m := newPagerModel(ctx, totalCount, model.PageSize, page).
			Add("Name", name).
			Add("Containers", containers)
		return ctx.Render("container/list", m)
	}

	c.Detail = func(ctx web.Context) error {
		id := ctx.P("id")
		container, err := docker.ContainerInspect(id)
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Container", container)
		return ctx.Render("container/detail", m)
	}

	c.Raw = func(ctx web.Context) error {
		id := ctx.P("id")
		container, raw, err := docker.ContainerInspectRaw(id)
		if err != nil {
			return err
		}

		j, err := misc.JSONIndent(raw)
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Container", container).Add("Raw", j)
		return ctx.Render("container/raw", m)
	}

	return
}

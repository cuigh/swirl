package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

type RegistryController struct {
	List   web.HandlerFunc `path:"/" name:"registry.list" authorize:"!" desc:"registry list page"`
	Create web.HandlerFunc `path:"/create" method:"post" name:"registry.create" authorize:"!" desc:"create registry"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"registry.delete" authorize:"!" desc:"delete registry"`
	Update web.HandlerFunc `path:"/update" method:"post" name:"registry.update" authorize:"!" desc:"update registry"`
}

func Registry() (c *RegistryController) {
	c = &RegistryController{}

	c.List = func(ctx web.Context) error {
		registries, err := biz.Registry.List()
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Registries", registries)
		return ctx.Render("registry/list", m)
	}

	c.Create = func(ctx web.Context) error {
		registry := &model.Registry{}
		err := ctx.Bind(registry)
		if err != nil {
			return err
		}
		err = biz.Registry.Create(registry, ctx.User())
		return ajaxResult(ctx, err)
	}

	c.Delete = func(ctx web.Context) error {
		id := ctx.F("id")
		err := biz.Registry.Delete(id, ctx.User())
		return ajaxResult(ctx, err)
	}

	c.Update = func(ctx web.Context) error {
		registry := &model.Registry{}
		err := ctx.Bind(registry)
		if err != nil {
			return err
		}
		err = biz.Registry.Update(registry, ctx.User())
		return ajaxResult(ctx, err)
	}

	return
}

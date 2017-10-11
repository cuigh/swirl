package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

// RegistryController is a controller of docker registry
type RegistryController struct {
	List   web.HandlerFunc `path:"/" name:"registry.list" authorize:"!" desc:"registry list page"`
	Create web.HandlerFunc `path:"/create" method:"post" name:"registry.create" authorize:"!" desc:"create registry"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"registry.delete" authorize:"!" desc:"delete registry"`
	Update web.HandlerFunc `path:"/update" method:"post" name:"registry.update" authorize:"!" desc:"update registry"`
}

// Registry creates an instance of RegistryController
func Registry() (c *RegistryController) {
	return &RegistryController{
		List:   registryList,
		Create: registryCreate,
		Delete: registryDelete,
		Update: registryUpdate,
	}
}

func registryList(ctx web.Context) error {
	registries, err := biz.Registry.List()
	if err != nil {
		return err
	}

	m := newModel(ctx).Add("Registries", registries)
	return ctx.Render("registry/list", m)
}

func registryCreate(ctx web.Context) error {
	registry := &model.Registry{}
	err := ctx.Bind(registry)
	if err != nil {
		return err
	}
	err = biz.Registry.Create(registry, ctx.User())
	return ajaxResult(ctx, err)
}

func registryDelete(ctx web.Context) error {
	id := ctx.F("id")
	err := biz.Registry.Delete(id, ctx.User())
	return ajaxResult(ctx, err)
}

func registryUpdate(ctx web.Context) error {
	registry := &model.Registry{}
	err := ctx.Bind(registry)
	if err != nil {
		return err
	}
	err = biz.Registry.Update(registry, ctx.User())
	return ajaxResult(ctx, err)
}

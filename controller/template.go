package controller

import (
	"encoding/json"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/model"
)

// TemplateController is a controller of service template
type TemplateController struct {
	List   web.HandlerFunc `path:"/" name:"template.list" authorize:"!" desc:"service template list page"`
	New    web.HandlerFunc `path:"/new" name:"template.new" authorize:"!" desc:"new service template page"`
	Create web.HandlerFunc `path:"/new" method:"post" name:"template.create" authorize:"!" desc:"create service template"`
	Edit   web.HandlerFunc `path:"/:id/edit" name:"template.edit" authorize:"!" desc:"edit service template page"`
	Update web.HandlerFunc `path:"/:id/edit" method:"post" name:"template.update" authorize:"!" desc:"update service template"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"template.delete" authorize:"!" desc:"delete service template"`
}

// Template creates an instance of TemplateController
func Template() (c *TemplateController) {
	return &TemplateController{
		List:   templateList,
		New:    templateNew,
		Create: templateCreate,
		Edit:   templateEdit,
		Update: templateUpdate,
		Delete: templateDelete,
	}
}

func templateList(ctx web.Context) error {
	args := &model.TemplateListArgs{}
	err := ctx.Bind(args)
	if err != nil {
		return err
	}

	args.PageSize = model.PageSize
	if args.PageIndex == 0 {
		args.PageIndex = 1
	}
	tpls, totalCount, err := biz.Template.List(args)
	if err != nil {
		return err
	}

	m := newPagerModel(ctx, totalCount, model.PageSize, args.PageIndex).
		Set("Name", args.Name).
		Set("Templates", tpls)
	return ctx.Render("service/template/list", m)
}

func templateNew(ctx web.Context) error {
	service := model.ServiceInfo{}
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
	m := newModel(ctx).Set("Action", "New").Set("Service", service).Set("Registries", registries).
		Set("Networks", networks).Set("CheckedNetworks", data.Set{}).
		Set("Secrets", secrets).Set("Configs", configs)
	return ctx.Render("service/template/edit", m)
}

func templateCreate(ctx web.Context) error {
	info := &model.ServiceInfo{}
	err := ctx.Bind(info)
	if err == nil {
		var (
			content []byte
			tpl     = &model.Template{Name: info.Name}
		)

		info.Normalize()
		info.Name = ""
		content, err = json.Marshal(info)
		if err != nil {
			return err
		}

		tpl.Content = string(content)
		err = biz.Template.Create(tpl, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func templateEdit(ctx web.Context) error {
	id := ctx.P("id")
	tpl, err := biz.Template.Get(id)
	if err != nil {
		return err
	} else if tpl == nil {
		return web.ErrNotFound
	}

	service := &model.ServiceInfo{}
	err = json.Unmarshal([]byte(tpl.Content), service)
	if err != nil {
		return err
	}
	service.Normalize()

	service.Name = tpl.Name
	if service.Registry != "" {
		var registry *model.Registry
		registry, err = biz.Registry.Get(service.Registry)
		if err != nil {
			return err
		}
		service.RegistryURL = registry.URL
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

	checkedNetworks := data.NewSet()
	checkedNetworks.AddSlice(service.Networks, func(i int) interface{} { return service.Networks[i] })

	m := newModel(ctx).Set("Action", "Edit").Set("Service", service).Set("Registries", registries).
		Set("Networks", networks).Set("CheckedNetworks", checkedNetworks).
		Set("Secrets", secrets).Set("Configs", configs)
	return ctx.Render("service/template/edit", m)
}

func templateUpdate(ctx web.Context) error {
	info := &model.ServiceInfo{}
	err := ctx.Bind(info)
	if err == nil {
		var (
			content []byte
			tpl     = &model.Template{
				ID:   ctx.P("id"),
				Name: info.Name,
			}
		)

		info.Normalize()
		info.Name = ""
		content, err = json.Marshal(info)
		if err != nil {
			return err
		}

		tpl.Content = string(content)
		err = biz.Template.Update(tpl, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func templateDelete(ctx web.Context) error {
	id := ctx.F("id")
	err := biz.Template.Delete(id, ctx.User())
	return ajaxResult(ctx, err)
}

package controller

import "github.com/cuigh/auxo/net/web"

type TemplateController struct {
	List web.HandlerFunc `path:"/" name:"template.list" authorize:"!" desc:"service template list page"`
}

func Template() (c *TemplateController) {
	c = &TemplateController{}

	c.List = func(ctx web.Context) error {
		m := newModel(ctx)
		return ctx.Render("service/template/list", m)
	}

	return
}

package controller

import (
	"strings"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/model"
)

type SecretController struct {
	List   web.HandlerFunc `path:"/" name:"secret.list" authorize:"!" desc:"secret list page"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"secret.delete" authorize:"!" desc:"delete secret"`
	New    web.HandlerFunc `path:"/new" name:"secret.new" authorize:"!" desc:"new secret page"`
	Create web.HandlerFunc `path:"/new" method:"post" name:"secret.create" authorize:"!" desc:"create secret"`
	Edit   web.HandlerFunc `path:"/:id/edit" name:"secret.edit" authorize:"!" desc:"edit secret page"`
	Update web.HandlerFunc `path:"/:id/update" method:"post" name:"secret.update" authorize:"!" desc:"update secret"`
}

func Secret() (c *SecretController) {
	c = &SecretController{}

	c.List = func(ctx web.Context) error {
		name := ctx.Q("name")
		page := cast.ToIntD(ctx.Q("page"), 1)
		secrets, totalCount, err := docker.SecretList(name, page, model.PageSize)
		if err != nil {
			return err
		}

		m := newPagerModel(ctx, totalCount, model.PageSize, page).
			Add("Name", name).
			Add("Secrets", secrets)
		return ctx.Render("secret/list", m)
	}

	c.Delete = func(ctx web.Context) error {
		ids := strings.Split(ctx.F("ids"), ",")
		for _, id := range ids {
			err := docker.SecretRemove(id)
			if err != nil {
				return ajaxResult(ctx, err)
			} else {
				// todo:
				biz.Event.CreateSecret(model.EventActionDelete, id, ctx.User())
			}
		}
		return ajaxSuccess(ctx, nil)
	}

	c.New = func(ctx web.Context) error {
		m := newModel(ctx)
		return ctx.Render("secret/new", m)
	}

	c.Create = func(ctx web.Context) error {
		v := &model.ConfigCreateInfo{}
		err := ctx.Bind(v)
		if err == nil {
			err = docker.SecretCreate(v)
			if err == nil {
				biz.Event.CreateSecret(model.EventActionCreate, v.Name, ctx.User())
			}
		}
		return ajaxResult(ctx, err)
	}

	c.Edit = func(ctx web.Context) error {
		id := ctx.P("id")
		secret, _, err := docker.SecretInspect(id)
		if err != nil {
			return err
		}
		m := newModel(ctx).Add("Secret", secret)
		return ctx.Render("secret/edit", m)
	}

	c.Update = func(ctx web.Context) error {
		v := &model.ConfigUpdateInfo{}
		err := ctx.Bind(v)
		if err == nil {
			err = docker.SecretUpdate(v)
			if err == nil {
				biz.Event.CreateSecret(model.EventActionUpdate, v.Name, ctx.User())
			}
		}
		return ajaxResult(ctx, err)
	}

	return
}

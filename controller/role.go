package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

// RoleController is a controller of user role
type RoleController struct {
	Index  web.HandlerFunc `path:"/" name:"role.list" authorize:"!" desc:"role list page"`
	New    web.HandlerFunc `path:"/new" name:"role.new" authorize:"!" desc:"new role page"`
	Create web.HandlerFunc `path:"/new" method:"post" name:"role.create" authorize:"!" desc:"create role"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"role.delete" authorize:"!" desc:"delete role"`
	Detail web.HandlerFunc `path:"/:id/detail" name:"role.detail" authorize:"!" desc:"role detail page"`
	Edit   web.HandlerFunc `path:"/:id/edit" name:"role.edit" authorize:"!" desc:"edit role page"`
	Update web.HandlerFunc `path:"/:id/update" method:"post" name:"role.update" authorize:"!" desc:"update role"`
}

// Role creates an instance of RoleController
func Role() (c *RoleController) {
	return &RoleController{
		Index:  roleIndex,
		New:    roleNew,
		Create: roleCreate,
		Delete: roleDelete,
		Detail: roleDetail,
		Edit:   roleEdit,
		Update: roleUpdate,
	}
}

func roleIndex(ctx web.Context) error {
	roles, err := biz.Role.List()
	if err != nil {
		return err
	}

	m := newModel(ctx).Add("Roles", roles)
	return ctx.Render("system/role/list", m)
}

func roleNew(ctx web.Context) error {
	m := newModel(ctx).Add("Perms", misc.Perms)
	return ctx.Render("system/role/new", m)
}

func roleCreate(ctx web.Context) error {
	role := &model.Role{}
	err := ctx.Bind(role)
	if err == nil {
		err = biz.Role.Create(role, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func roleDelete(ctx web.Context) error {
	id := ctx.F("id")
	err := biz.Role.Delete(id, ctx.User())
	return ajaxResult(ctx, err)
}

func roleDetail(ctx web.Context) error {
	id := ctx.P("id")
	role, err := biz.Role.Get(id)
	if err != nil {
		return err
	}
	if role == nil {
		return web.ErrNotFound
	}

	perms := make(map[string]struct{})
	for _, p := range role.Perms {
		perms[p] = model.Placeholder
	}
	m := newModel(ctx).Add("Role", role).Add("Perms", misc.Perms).Add("CheckedPerms", perms)
	return ctx.Render("system/role/detail", m)
}

func roleEdit(ctx web.Context) error {
	id := ctx.P("id")
	role, err := biz.Role.Get(id)
	if err != nil {
		return err
	}
	if role == nil {
		return web.ErrNotFound
	}

	perms := make(map[string]struct{})
	for _, p := range role.Perms {
		perms[p] = model.Placeholder
	}
	m := newModel(ctx).Add("Role", role).Add("Perms", misc.Perms).Add("CheckedPerms", perms)
	return ctx.Render("system/role/edit", m)
}

func roleUpdate(ctx web.Context) error {
	role := &model.Role{}
	err := ctx.Bind(role)
	if err == nil {
		role.ID = ctx.P("id")
		err = biz.Role.Update(role, ctx.User())
	}
	return ajaxResult(ctx, err)
}

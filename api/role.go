package api

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
)

// RoleHandler encapsulates role related handlers.
type RoleHandler struct {
	Find   web.HandlerFunc `path:"/find" auth:"role.view" desc:"find role by id"`
	Search web.HandlerFunc `path:"/search" auth:"role.view" desc:"search roles"`
	Save   web.HandlerFunc `path:"/save" method:"post" auth:"role.edit" desc:"create or update role"`
	Delete web.HandlerFunc `path:"/delete" method:"post" auth:"role.delete" desc:"delete role"`
}

// NewRole creates an instance of RoleHandler
func NewRole(b biz.RoleBiz) *RoleHandler {
	return &RoleHandler{
		Search: roleSearch(b),
		Find:   roleFind(b),
		Delete: roleDelete(b),
		Save:   roleSave(b),
	}
}

func roleSearch(b biz.RoleBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		name := ctx.Query("name")
		roles, err := b.Search(name)
		if err != nil {
			return err
		}
		return success(ctx, roles)
	}
}

func roleFind(b biz.RoleBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		id := ctx.Query("id")
		role, err := b.Find(id)
		if err != nil {
			return err
		}
		return success(ctx, role)
	}
}

func roleDelete(b biz.RoleBiz) web.HandlerFunc {
	type Args struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Delete(args.ID, args.Name, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func roleSave(b biz.RoleBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		r := &biz.Role{}
		err := ctx.Bind(r, true)
		if err == nil {
			if r.ID == "" {
				err = b.Create(r, ctx.User())
			} else {
				err = b.Update(r, ctx.User())
			}
		}
		return ajax(ctx, err)
	}
}

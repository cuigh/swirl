package api

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
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
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		name := c.Query("name")
		roles, err := b.Search(ctx, name)
		if err != nil {
			return err
		}
		return success(c, roles)
	}
}

func roleFind(b biz.RoleBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		id := c.Query("id")
		role, err := b.Find(ctx, id)
		if err != nil {
			return err
		}
		return success(c, role)
	}
}

func roleDelete(b biz.RoleBiz) web.HandlerFunc {
	type Args struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Delete(ctx, args.ID, args.Name, c.User())
		}
		return ajax(c, err)
	}
}

func roleSave(b biz.RoleBiz) web.HandlerFunc {
	return func(c web.Context) error {
		r := &dao.Role{}
		err := c.Bind(r, true)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			if r.ID == "" {
				err = b.Create(ctx, r, c.User())
			} else {
				err = b.Update(ctx, r, c.User())
			}
		}
		return ajax(c, err)
	}
}

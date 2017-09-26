package controller

import (
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

type UserController struct {
	Index   web.HandlerFunc `path:"/" name:"user.list" authorize:"!" desc:"user list page"`
	New     web.HandlerFunc `path:"/new" name:"user.new" authorize:"!" desc:"new user page"`
	Create  web.HandlerFunc `path:"/new" method:"post" name:"user.create" authorize:"!" desc:"create user"`
	Detail  web.HandlerFunc `path:"/:name/detail" name:"user.detail" authorize:"!" desc:"user detail page"`
	Edit    web.HandlerFunc `path:"/:name/edit" name:"user.edit" authorize:"!" desc:"edit user page"`
	Update  web.HandlerFunc `path:"/:name/update" method:"post" name:"user.update" authorize:"!" desc:"update user"`
	Block   web.HandlerFunc `path:"/block" method:"post" name:"user.block" authorize:"!" desc:"block user"`
	Unblock web.HandlerFunc `path:"/unblock" method:"post" name:"user.unblock" authorize:"!" desc:"unblock user"`
	Delete  web.HandlerFunc `path:"/delete" method:"post" name:"user.delete" authorize:"!" desc:"delete user"`
}

func User() (c *UserController) {
	c = &UserController{}

	c.Index = func(ctx web.Context) error {
		args := &model.UserListArgs{}
		err := ctx.Bind(args)
		if err != nil {
			return err
		}
		args.PageSize = model.PageSize
		if args.PageIndex == 0 {
			args.PageIndex = 1
		}

		users, totalCount, err := biz.User.List(args)
		if err != nil {
			return err
		}

		m := newPagerModel(ctx, totalCount, model.PageSize, args.PageIndex).
			Add("Query", args.Query).
			Add("Filter", args.Filter).
			Add("Users", users)
		return ctx.Render("system/user/list", m)
	}

	c.New = func(ctx web.Context) error {
		roles, err := biz.Role.List()
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Roles", roles)
		return ctx.Render("system/user/new", m)
	}

	c.Create = func(ctx web.Context) error {
		user := &model.User{}
		err := ctx.Bind(user)
		if err == nil {
			user.Type = model.UserTypeInternal
			err = biz.User.Create(user, ctx.User())
		}
		return ajaxResult(ctx, err)
	}

	c.Detail = func(ctx web.Context) error {
		name := ctx.P("name")
		user, err := biz.User.GetByName(name)
		if err != nil {
			return err
		}
		if user == nil {
			return web.ErrNotFound
		}

		var (
			roles map[string]string
			role  *model.Role
		)
		if len(user.Roles) > 0 {
			roles = map[string]string{}
			for _, id := range user.Roles {
				role, err = biz.Role.Get(id)
				if err != nil {
					return err
				} else if role != nil {
					roles[id] = role.Name
				} else {
					log.Get("user").Warnf("Role %v is invalid", id)
				}
			}
		}

		m := newModel(ctx).Add("User", user).Add("Roles", roles)
		return ctx.Render("system/user/detail", m)
	}

	c.Edit = func(ctx web.Context) error {
		name := ctx.P("name")
		user, err := biz.User.GetByName(name)
		if err != nil {
			return err
		}
		if user == nil {
			return web.ErrNotFound
		}

		roles, err := biz.Role.List()
		if err != nil {
			return err
		}

		userRoles := make(map[string]struct{})
		for _, id := range user.Roles {
			userRoles[id] = model.Placeholder
		}
		m := newModel(ctx).Add("User", user).Add("Roles", roles).Add("UserRoles", userRoles)
		return ctx.Render("system/user/edit", m)
	}

	c.Update = func(ctx web.Context) error {
		user := &model.User{}
		err := ctx.Bind(user)
		if err == nil {
			err = biz.User.Update(user, ctx.User())
		}
		return ajaxResult(ctx, err)
	}

	c.Block = func(ctx web.Context) error {
		id := ctx.F("id")
		err := biz.User.Block(id)
		return ajaxResult(ctx, err)
	}

	c.Unblock = func(ctx web.Context) error {
		id := ctx.F("id")
		err := biz.User.Unblock(id)
		return ajaxResult(ctx, err)
	}

	c.Delete = func(ctx web.Context) error {
		id := ctx.F("id")
		err := biz.User.Delete(id)
		return ajaxResult(ctx, err)
	}

	return
}

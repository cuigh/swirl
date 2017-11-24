package controller

import (
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

// UserController is a controller of user
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
	Search  web.HandlerFunc `path:"/search" method:"post" name:"user.search" authorize:"?" desc:"search users"`
}

// User creates an instance of UserController
func User() (c *UserController) {
	return &UserController{
		Index:   userIndex,
		New:     userNew,
		Create:  userCreate,
		Detail:  userDetail,
		Edit:    userEdit,
		Update:  userUpdate,
		Block:   userBlock,
		Unblock: userUnblock,
		Delete:  userDelete,
		Search:  userSearch,
	}
}

func userIndex(ctx web.Context) error {
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
		Set("Query", args.Query).
		Set("Filter", args.Filter).
		Set("Users", users)
	return ctx.Render("system/user/list", m)
}

func userNew(ctx web.Context) error {
	roles, err := biz.Role.List()
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Roles", roles)
	return ctx.Render("system/user/new", m)
}

func userCreate(ctx web.Context) error {
	user := &model.User{}
	err := ctx.Bind(user, true)
	if err == nil {
		user.Type = model.UserTypeInternal
		err = biz.User.Create(user, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func userDetail(ctx web.Context) error {
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

	m := newModel(ctx).Set("User", user).Set("Roles", roles)
	return ctx.Render("system/user/detail", m)
}

func userEdit(ctx web.Context) error {
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
	m := newModel(ctx).Set("User", user).Set("Roles", roles).Set("UserRoles", userRoles)
	return ctx.Render("system/user/edit", m)
}

func userUpdate(ctx web.Context) error {
	user := &model.User{}
	err := ctx.Bind(user)
	if err == nil {
		err = biz.User.Update(user, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func userBlock(ctx web.Context) error {
	id := ctx.F("id")
	err := biz.User.Block(id)
	return ajaxResult(ctx, err)
}

func userUnblock(ctx web.Context) error {
	id := ctx.F("id")
	err := biz.User.Unblock(id)
	return ajaxResult(ctx, err)
}

func userDelete(ctx web.Context) error {
	id := ctx.F("id")
	err := biz.User.Delete(id)
	return ajaxResult(ctx, err)
}

func userSearch(ctx web.Context) error {
	query := ctx.F("query")
	args := &model.UserListArgs{
		Query:     query,
		PageIndex: 1,
		PageSize:  10,
	}
	users, _, err := biz.User.List(args)
	if err != nil {
		return err
	}

	type User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	list := make([]User, len(users))
	for i, user := range users {
		list[i] = User{
			ID:   user.ID,
			Name: user.Name,
		}
	}
	return ctx.JSON(list)
}

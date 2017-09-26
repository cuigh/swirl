package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

type ProfileController struct {
	Index          web.HandlerFunc `path:"/" name:"profile.info" authorize:"?" desc:"profile info page"`
	ModifyInfo     web.HandlerFunc `path:"/" method:"post" name:"profile.info.modify" authorize:"?" desc:"modify info"`
	Password       web.HandlerFunc `path:"/password" name:"profile.password" authorize:"?" desc:"profile password page"`
	ModifyPassword web.HandlerFunc `path:"/password" method:"post" name:"profile.password.modify" authorize:"?" desc:"modify password"`
}

func Profile() (c *ProfileController) {
	c = &ProfileController{}

	c.Index = func(ctx web.Context) error {
		user, err := biz.User.GetByID(ctx.User().ID())
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("User", user)
		return ctx.Render("profile/index", m)
	}

	c.ModifyInfo = func(ctx web.Context) error {
		user := &model.User{}
		err := ctx.Bind(user)
		if err == nil {
			user.ID = ctx.User().ID()
			err = biz.User.UpdateInfo(user)
		}
		return ajaxResult(ctx, err)
	}

	c.Password = func(ctx web.Context) error {
		m := newModel(ctx)
		return ctx.Render("profile/password", m)
	}

	c.ModifyPassword = func(ctx web.Context) error {
		old_pwd := ctx.F("password_old")
		new_pwd := ctx.F("password")
		err := biz.User.UpdatePassword(ctx.User().ID(), old_pwd, new_pwd)
		return ajaxResult(ctx, err)
	}

	return
}

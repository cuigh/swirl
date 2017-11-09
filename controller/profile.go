package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

// ProfileController is a controller of user profile
type ProfileController struct {
	Index          web.HandlerFunc `path:"/" name:"profile.info" authorize:"?" desc:"profile info page"`
	ModifyInfo     web.HandlerFunc `path:"/" method:"post" name:"profile.info.modify" authorize:"?" desc:"modify info"`
	Password       web.HandlerFunc `path:"/password" name:"profile.password" authorize:"?" desc:"profile password page"`
	ModifyPassword web.HandlerFunc `path:"/password" method:"post" name:"profile.password.modify" authorize:"?" desc:"modify password"`
}

// Profile creates an instance of ProfileController
func Profile() (c *ProfileController) {
	return &ProfileController{
		Index:          profileIndex,
		ModifyInfo:     profileModifyInfo,
		Password:       profilePassword,
		ModifyPassword: profileModifyPassword,
	}
}

func profileIndex(ctx web.Context) error {
	user, err := biz.User.GetByID(ctx.User().ID())
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("User", user)
	return ctx.Render("profile/index", m)
}

func profileModifyInfo(ctx web.Context) error {
	user := &model.User{}
	err := ctx.Bind(user)
	if err == nil {
		user.ID = ctx.User().ID()
		err = biz.User.UpdateInfo(user)
	}
	return ajaxResult(ctx, err)
}

func profilePassword(ctx web.Context) error {
	user, err := biz.User.GetByID(ctx.User().ID())
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("User", user)
	return ctx.Render("profile/password", m)
}

func profileModifyPassword(ctx web.Context) error {
	oldPwd := ctx.F("password_old")
	newPwd := ctx.F("password")
	err := biz.User.UpdatePassword(ctx.User().ID(), oldPwd, newPwd)
	return ajaxResult(ctx, err)
}

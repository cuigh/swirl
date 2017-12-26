package controller

import (
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/model"
)

// HomeController is a basic controller of site
type HomeController struct {
	Index    web.HandlerFunc `path:"/" name:"index" authorize:"?" desc:"index page"`
	Login    web.HandlerFunc `path:"/login" name:"login.view" authorize:"*" desc:"sign in page"`
	InitGet  web.HandlerFunc `path:"/init" name:"init.view" authorize:"*" desc:"initialize page"`
	InitPost web.HandlerFunc `path:"/init" name:"init" method:"post" authorize:"*" desc:"initialize system"`
	Error403 web.HandlerFunc `path:"/403" name:"403" authorize:"?" desc:"403 page"`
	Error404 web.HandlerFunc `path:"/404" name:"404" authorize:"*" desc:"404 page"`
}

// Home creates an instance of HomeController
func Home() (c *HomeController) {
	return &HomeController{
		Index:    homeIndex,
		Login:    homeLogin,
		InitGet:  homeInitGet,
		InitPost: homeInitPost,
		Error403: homeError403,
		Error404: homeError404,
	}
}

func homeIndex(ctx web.Context) (err error) {
	var (
		count int
		m     = newModel(ctx)
	)

	if count, err = docker.NodeCount(); err != nil {
		return
	}
	m.Set("NodeCount", count)

	if count, err = docker.NetworkCount(); err != nil {
		return
	}
	m.Set("NetworkCount", count)

	if count, err = docker.ServiceCount(); err != nil {
		return
	}
	m.Set("ServiceCount", count)

	if count, err = docker.StackCount(); err != nil {
		return
	}
	m.Set("StackCount", count)

	return ctx.Render("index", m)
}

func homeLogin(ctx web.Context) error {
	count, err := biz.User.Count()
	if err != nil {
		return err
	} else if count == 0 {
		return ctx.Redirect("init")
	}
	if ctx.User() != nil {
		u := ctx.Q("from")
		if u == "" {
			u = "/"
		}
		return ctx.Redirect(u)
	}
	return ctx.Render("login", nil)
}

func homeInitGet(ctx web.Context) error {
	count, err := biz.User.Count()
	if err != nil {
		return err
	} else if count > 0 {
		return ctx.Redirect("login")
	}
	return ctx.Render("init", nil)
}

func homeInitPost(ctx web.Context) error {
	count, err := biz.User.Count()
	if err != nil {
		return err
	} else if count > 0 {
		return errors.New("Swirl was already initialized")
	}

	user := &model.User{}
	err = ctx.Bind(user)
	if err != nil {
		return err
	}

	user.Admin = true
	user.Type = model.UserTypeInternal
	err = biz.User.Create(user, nil)
	return ajaxResult(ctx, err)
}

func homeError403(ctx web.Context) error {
	return ctx.Render("403", nil)
}

func homeError404(ctx web.Context) error {
	return ctx.Render("404", nil)
}

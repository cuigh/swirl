package controller

import (
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/model"
)

type HomeController struct {
	Index    web.HandlerFunc `path:"/" name:"index" authorize:"?" desc:"index page"`
	Error404 web.HandlerFunc `path:"/404" name:"404" authorize:"*" desc:"404 page"`
	Login    web.HandlerFunc `path:"/login" name:"login" authorize:"*" desc:"sign in page"`
	InitGet  web.HandlerFunc `path:"/init" name:"init" authorize:"*" desc:"initialize page"`
	InitPost web.HandlerFunc `path:"/init" method:"post" name:"init" authorize:"*" desc:"initialize system"`
}

func Home() (c *HomeController) {
	c = &HomeController{}

	c.Index = func(ctx web.Context) (err error) {
		var (
			count int
			m     = newModel(ctx)
		)

		if count, err = docker.NodeCount(); err != nil {
			return
		}
		m.Add("NodeCount", count)

		if count, err = docker.NetworkCount(); err != nil {
			return
		}
		m.Add("NetworkCount", count)

		if count, err = docker.ServiceCount(); err != nil {
			return
		}
		m.Add("ServiceCount", count)

		if count, err = docker.StackCount(); err != nil {
			return
		}
		m.Add("StackCount", count)

		return ctx.Render("index", m)
	}

	c.Login = func(ctx web.Context) error {
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

	c.InitGet = func(ctx web.Context) error {
		count, err := biz.User.Count()
		if err != nil {
			return err
		} else if count > 0 {
			return ctx.Redirect("login")
		}
		return ctx.Render("init", nil)
	}

	c.InitPost = func(ctx web.Context) error {
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

	c.Error404 = func(ctx web.Context) error {
		return ctx.Render("404", nil)
	}
	return
}

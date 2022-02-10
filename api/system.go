package api

import (
	"runtime"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/docker"
	"github.com/cuigh/swirl/misc"
)

// SystemHandler encapsulates system related handlers.
type SystemHandler struct {
	CheckState  web.HandlerFunc `path:"/check-state" auth:"*" desc:"check system state"`
	CreateAdmin web.HandlerFunc `path:"/create-admin" method:"post" auth:"*" desc:"initialize administrator account"`
	Version     web.HandlerFunc `path:"/version" auth:"*" desc:"fetch app version"`
	Summarize   web.HandlerFunc `path:"/summarize" auth:"?" desc:"fetch statistics data"`
}

// NewSystem creates an instance of SystemHandler
func NewSystem(d *docker.Docker, b biz.SystemBiz, ub biz.UserBiz) *SystemHandler {
	return &SystemHandler{
		CheckState:  systemCheckState(b),
		CreateAdmin: systemCreateAdmin(ub),
		Version:     systemVersion,
		Summarize:   systemSummarize(d),
	}
}

func systemCheckState(b biz.SystemBiz) web.HandlerFunc {
	return func(c web.Context) (err error) {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		state, err := b.CheckState(ctx)
		if err != nil {
			return err
		}
		return success(c, state)
	}
}

func systemVersion(c web.Context) (err error) {
	return success(c, data.Map{
		"version":   app.Version,
		"goVersion": runtime.Version(),
	})
}

func systemSummarize(d *docker.Docker) web.HandlerFunc {
	return func(c web.Context) (err error) {
		summary := struct {
			NodeCount    int `json:"nodeCount"`
			NetworkCount int `json:"networkCount"`
			ServiceCount int `json:"serviceCount"`
			StackCount   int `json:"stackCount"`
		}{}

		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		if summary.NodeCount, err = d.NodeCount(ctx); err != nil {
			return
		}
		if summary.NetworkCount, err = d.NetworkCount(ctx); err != nil {
			return
		}
		if summary.ServiceCount, err = d.ServiceCount(ctx); err != nil {
			return
		}
		if summary.StackCount, err = d.StackCount(ctx); err != nil {
			return
		}

		return success(c, summary)
	}
}

func systemCreateAdmin(ub biz.UserBiz) web.HandlerFunc {
	return func(c web.Context) (err error) {
		args := &struct {
			Password string `json:"password"`
			*dao.User
		}{}
		if err = c.Bind(args, true); err != nil {
			return err
		}

		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		var count int
		if count, err = ub.Count(ctx); err == nil && count > 0 {
			return errors.Coded(misc.ErrSystemInitialized, "system was already initialized")
		}

		user := args.User
		user.Password = args.Password
		user.Admin = true
		user.Type = biz.UserTypeInternal
		_, err = ub.Create(ctx, user, nil)
		return ajax(c, err)
	}
}

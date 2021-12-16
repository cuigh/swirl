package api

import (
	"context"
	"runtime"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/docker"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

//var ErrSystemInitialized = errors.New("system was already initialized")

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
		state, err := b.CheckState()
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

		if summary.NodeCount, err = d.NodeCount(context.TODO()); err != nil {
			return
		}
		if summary.NetworkCount, err = d.NetworkCount(context.TODO()); err != nil {
			return
		}
		if summary.ServiceCount, err = d.ServiceCount(context.TODO()); err != nil {
			return
		}
		if summary.StackCount, err = d.StackCount(context.TODO()); err != nil {
			return
		}

		return success(c, summary)
	}
}

func systemCreateAdmin(ub biz.UserBiz) web.HandlerFunc {
	return func(c web.Context) (err error) {
		user := &model.User{}
		if err = c.Bind(user, true); err != nil {
			return err
		}

		var count int
		if count, err = ub.Count(); err == nil && count > 0 {
			return errors.Coded(misc.ErrSystemInitialized, "system was already initialized")
		}

		user.Admin = true
		user.Type = biz.UserTypeInternal
		_, err = ub.Create(user, nil)
		return ajax(c, err)
	}
}

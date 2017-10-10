package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/misc"
)

type TaskController struct {
	Detail web.HandlerFunc `path:"/:id/detail" name:"task.detail" authorize:"!" desc:"task detail page"`
	Raw    web.HandlerFunc `path:"/:id/raw" name:"task.raw" authorize:"!" desc:"task raw page"`
}

func Task() (c *TaskController) {
	c = &TaskController{}

	c.Detail = func(ctx web.Context) error {
		id := ctx.P("id")
		task, _, err := docker.TaskInspect(id)
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Task", task)
		return ctx.Render("task/detail", m)
	}

	c.Raw = func(ctx web.Context) error {
		id := ctx.P("id")
		task, raw, err := docker.TaskInspect(id)
		if err != nil {
			return err
		}

		j, err := misc.JSONIndent(raw)
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Task", task).Add("Raw", j)
		return ctx.Render("task/raw", m)
	}

	return
}

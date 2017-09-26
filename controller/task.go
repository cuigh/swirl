package controller

import (
	"bytes"
	"encoding/json"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz/docker"
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

		buf := &bytes.Buffer{}
		err = json.Indent(buf, raw, "", "    ")
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Task", task).Add("Raw", string(buf.Bytes()))
		return ctx.Render("task/raw", m)
	}

	return
}

package controller

import (
	"strings"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

// ContainerController is a controller of docker container
type ContainerController struct {
	List   web.HandlerFunc `path:"/" name:"container.list" authorize:"!" desc:"container list page"`
	Detail web.HandlerFunc `path:"/:id/detail" name:"container.detail" authorize:"!" desc:"container detail page"`
	Raw    web.HandlerFunc `path:"/:id/raw" name:"container.raw" authorize:"!" desc:"container raw page"`
	Logs   web.HandlerFunc `path:"/:id/logs" name:"container.logs" authorize:"!" desc:"container logs page"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"container.delete" authorize:"!" desc:"delete container"`
}

// Container creates an instance of ContainerController
func Container() (c *ContainerController) {
	return &ContainerController{
		List:   containerList,
		Detail: containerDetail,
		Raw:    containerRaw,
		Logs:   containerLogs,
		Delete: containerDelete,
	}
}

func containerList(ctx web.Context) error {
	args := &model.ContainerListArgs{}
	err := ctx.Bind(args)
	if err != nil {
		return err
	}
	args.PageSize = model.PageSize
	if args.PageIndex == 0 {
		args.PageIndex = 1
	}

	containers, totalCount, err := docker.ContainerList(args)
	if err != nil {
		return err
	}

	m := newPagerModel(ctx, totalCount, model.PageSize, args.PageIndex).
		Add("Name", args.Name).
		Add("Filter", args.Filter).
		Add("Containers", containers)
	return ctx.Render("container/list", m)
}

func containerDetail(ctx web.Context) error {
	id := ctx.P("id")
	container, err := docker.ContainerInspect(id)
	if err != nil {
		return err
	}

	m := newModel(ctx).Add("Container", container)
	return ctx.Render("container/detail", m)
}

func containerRaw(ctx web.Context) error {
	id := ctx.P("id")
	container, raw, err := docker.ContainerInspectRaw(id)
	if err != nil {
		return err
	}

	j, err := misc.JSONIndent(raw)
	if err != nil {
		return err
	}

	m := newModel(ctx).Add("Container", container).Add("Raw", j)
	return ctx.Render("container/raw", m)
}

func containerLogs(ctx web.Context) error {
	id := ctx.P("id")
	container, _, err := docker.ContainerInspectRaw(id)
	if err != nil {
		return err
	}

	line := cast.ToInt(ctx.Q("line"), 500)
	timestamps := cast.ToBool(ctx.Q("timestamps"), false)
	stdout, stderr, err := docker.ContainerLogs(id, line, timestamps)
	if err != nil {
		return err
	}

	m := newModel(ctx).Add("Container", container).Add("Line", line).Add("Timestamps", timestamps).
		Add("Stdout", stdout.String()).Add("Stderr", stderr.String())
	return ctx.Render("container/logs", m)
}

func containerDelete(ctx web.Context) error {
	ids := strings.Split(ctx.F("ids"), ",")
	for _, id := range ids {
		if err := docker.ContainerRemove(id); err != nil {
			return ajaxResult(ctx, err)
		}
	}
	return ajaxSuccess(ctx, nil)
}

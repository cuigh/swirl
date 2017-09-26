package controller

import (
	"bytes"
	"encoding/json"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/model"
)

type NodeController struct {
	List   web.HandlerFunc `path:"/" name:"node.list" authorize:"!" desc:"node list page"`
	Detail web.HandlerFunc `path:"/:id/detail" name:"node.detail" authorize:"!" desc:"node detail page"`
	Raw    web.HandlerFunc `path:"/:id/raw" name:"node.raw" authorize:"!" desc:"node raw page"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"node.delete" authorize:"!" desc:"delete node"`
	Edit   web.HandlerFunc `path:"/:id/edit" name:"node.edit" authorize:"!" desc:"node edit page"`
	Update web.HandlerFunc `path:"/:id/update" method:"post" name:"node.update" authorize:"!" desc:"update node"`
}

func Node() (c *NodeController) {
	c = &NodeController{}

	c.List = func(ctx web.Context) error {
		nodes, err := docker.NodeList()
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Nodes", nodes)
		return ctx.Render("node/list", m)
	}

	c.Delete = func(ctx web.Context) error {
		id := ctx.F("id")
		err := docker.NodeRemove(id)
		return ajaxResult(ctx, err)
	}

	c.Detail = func(ctx web.Context) error {
		id := ctx.P("id")
		node, _, err := docker.NodeInspect(id)
		if err != nil {
			return err
		}

		tasks, err := docker.TaskList("", id)
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Node", node).Add("Tasks", tasks)
		return ctx.Render("node/detail", m)
	}

	c.Raw = func(ctx web.Context) error {
		id := ctx.P("id")
		node, raw, err := docker.NodeInspect(id)
		if err != nil {
			return err
		}

		buf := &bytes.Buffer{}
		err = json.Indent(buf, raw, "", "    ")
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("ID", id).Add("Node", node).Add("Raw", string(buf.Bytes()))
		return ctx.Render("node/raw", m)
	}

	c.Edit = func(ctx web.Context) error {
		id := ctx.P("id")
		node, _, err := docker.NodeInspect(id)
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Node", node)
		return ctx.Render("node/edit", m)
	}

	c.Update = func(ctx web.Context) error {
		id := ctx.P("id")
		info := &model.NodeUpdateInfo{}
		err := ctx.Bind(info)
		if err == nil {
			err = docker.NodeUpdate(id, info)
		}
		return ajaxResult(ctx, err)
	}

	return
}

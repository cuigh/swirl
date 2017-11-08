package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

// NodeController is a controller of swarm node
type NodeController struct {
	List   web.HandlerFunc `path:"/" name:"node.list" authorize:"!" desc:"node list page"`
	Detail web.HandlerFunc `path:"/:id/detail" name:"node.detail" authorize:"!" desc:"node detail page"`
	Raw    web.HandlerFunc `path:"/:id/raw" name:"node.raw" authorize:"!" desc:"node raw page"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"node.delete" authorize:"!" desc:"delete node"`
	Edit   web.HandlerFunc `path:"/:id/edit" name:"node.edit" authorize:"!" desc:"node edit page"`
	Update web.HandlerFunc `path:"/:id/update" method:"post" name:"node.update" authorize:"!" desc:"update node"`
}

// Node creates an instance of NodeController
func Node() (c *NodeController) {
	return &NodeController{
		List:   nodeList,
		Detail: nodeDetail,
		Raw:    nodeRaw,
		Delete: nodeDelete,
		Edit:   nodeEdit,
		Update: nodeUpdate,
	}
}

func nodeList(ctx web.Context) error {
	nodes, err := docker.NodeList()
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Nodes", nodes)
	return ctx.Render("node/list", m)
}

func nodeDelete(ctx web.Context) error {
	id := ctx.F("id")
	err := docker.NodeRemove(id)
	return ajaxResult(ctx, err)
}

func nodeDetail(ctx web.Context) error {
	id := ctx.P("id")
	node, _, err := docker.NodeInspect(id)
	if err != nil {
		return err
	}

	tasks, _, err := docker.TaskList(&model.TaskListArgs{Node: id})
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Node", node).Set("Tasks", tasks)
	return ctx.Render("node/detail", m)
}

func nodeRaw(ctx web.Context) error {
	id := ctx.P("id")
	node, raw, err := docker.NodeInspect(id)
	if err != nil {
		return err
	}

	j, err := misc.JSONIndent(raw)
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("ID", id).Set("Node", node).Set("Raw", j)
	return ctx.Render("node/raw", m)
}

func nodeEdit(ctx web.Context) error {
	id := ctx.P("id")
	node, _, err := docker.NodeInspect(id)
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Node", node)
	return ctx.Render("node/edit", m)
}

func nodeUpdate(ctx web.Context) error {
	id := ctx.P("id")
	info := &model.NodeUpdateInfo{}
	err := ctx.Bind(info)
	if err == nil {
		err = docker.NodeUpdate(id, info)
	}
	return ajaxResult(ctx, err)
}

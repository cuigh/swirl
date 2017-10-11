package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

// NetworkController is a controller of docker network
type NetworkController struct {
	List       web.HandlerFunc `path:"/" name:"network.list" authorize:"!" desc:"network list page"`
	New        web.HandlerFunc `path:"/new" name:"network.new" authorize:"!" desc:"new network page"`
	Create     web.HandlerFunc `path:"/create" method:"post" name:"network.create" authorize:"!" desc:"create network"`
	Delete     web.HandlerFunc `path:"/delete" method:"post" name:"network.delete" authorize:"!" desc:"delete network"`
	Disconnect web.HandlerFunc `path:"/:name/disconnect" method:"post" name:"network.disconnect" authorize:"!" desc:"disconnect network"`
	Detail     web.HandlerFunc `path:"/:name/detail" name:"network.detail" authorize:"!" desc:"network detail page"`
	Raw        web.HandlerFunc `path:"/:name/raw" name:"network.raw" authorize:"!" desc:"network raw page"`
}

// Network creates a NetworkController instance.
func Network() (c *NetworkController) {
	return &NetworkController{
		List:       networkList,
		New:        networkNew,
		Create:     networkCreate,
		Delete:     networkDelete,
		Disconnect: networkDisconnect,
		Detail:     networkDetail,
		Raw:        networkRaw,
	}
}

func networkList(ctx web.Context) error {
	networks, err := docker.NetworkList()
	if err != nil {
		return err
	}

	m := newModel(ctx).Add("Networks", networks)
	return ctx.Render("network/list", m)
}

func networkNew(ctx web.Context) error {
	m := newModel(ctx)
	return ctx.Render("/network/new", m)
}

func networkCreate(ctx web.Context) error {
	info := &model.NetworkCreateInfo{}
	err := ctx.Bind(info)
	if err != nil {
		return err
	}
	err = docker.NetworkCreate(info)
	if err == nil {
		biz.Event.CreateNetwork(model.EventActionCreate, info.Name, info.Name, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func networkDelete(ctx web.Context) error {
	name := ctx.F("name")
	err := docker.NetworkRemove(name)
	if err == nil {
		biz.Event.CreateNetwork(model.EventActionDelete, name, name, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func networkDisconnect(ctx web.Context) error {
	name := ctx.P("name")
	container := ctx.F("container")
	err := docker.NetworkDisconnect(name, container)
	if err == nil {
		biz.Event.CreateNetwork(model.EventActionDisconnect, name, name+" <-> "+container, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func networkDetail(ctx web.Context) error {
	name := ctx.P("name")
	network, err := docker.NetworkInspect(name)
	if err != nil {
		return err
	}
	m := newModel(ctx).Add("Network", network)
	return ctx.Render("network/detail", m)
}

func networkRaw(ctx web.Context) error {
	name := ctx.P("name")
	raw, err := docker.NetworkInspectRaw(name)
	if err != nil {
		return err
	}

	j, err := misc.JSONIndent(raw)
	if err != nil {
		return err
	}

	m := newModel(ctx).Add("Network", name).Add("Raw", j)
	return ctx.Render("network/raw", m)
}

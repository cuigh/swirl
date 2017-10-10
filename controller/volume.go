package controller

import (
	"strings"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

type VolumeController struct {
	List   web.HandlerFunc `path:"/" name:"volume.list" authorize:"!" desc:"volume list page"`
	New    web.HandlerFunc `path:"/new" name:"volume.new" authorize:"!" desc:"new volume page"`
	Create web.HandlerFunc `path:"/create" method:"post" name:"volume.create" authorize:"!" desc:"create volume"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"volume.delete" authorize:"!" desc:"delete volume"`
	Prune  web.HandlerFunc `path:"/prune" method:"post" name:"volume.prune" authorize:"!" desc:"prune volume"`
	Detail web.HandlerFunc `path:"/:name/detail" name:"volume.detail" authorize:"!" desc:"volume detail page"`
	Raw    web.HandlerFunc `path:"/:name/raw" name:"volume.raw" authorize:"!" desc:"volume raw page"`
}

func Volume() (c *VolumeController) {
	c = &VolumeController{}

	c.List = func(ctx web.Context) error {
		//name := ctx.Q("name")
		page := cast.ToIntD(ctx.Q("page"), 1)
		volumes, totalCount, err := docker.VolumeList(page, model.PageSize)
		if err != nil {
			return err
		}

		m := newPagerModel(ctx, totalCount, model.PageSize, page).
			//Add("Name", name).
			Add("Volumes", volumes)
		return ctx.Render("volume/list", m)
	}

	c.New = func(ctx web.Context) error {
		m := newModel(ctx)
		return ctx.Render("/volume/new", m)
	}

	c.Create = func(ctx web.Context) error {
		info := &model.VolumeCreateInfo{}
		err := ctx.Bind(info)
		if err != nil {
			return err
		}
		err = docker.VolumeCreate(info)
		if err == nil {
			biz.Event.CreateVolume(model.EventActionCreate, info.Name, ctx.User())
		}
		return ajaxResult(ctx, err)
	}

	c.Delete = func(ctx web.Context) error {
		names := strings.Split(ctx.F("names"), ",")
		for _, name := range names {
			if err := docker.VolumeRemove(name); err != nil {
				return ajaxResult(ctx, err)
			} else {
				biz.Event.CreateVolume(model.EventActionDelete, name, ctx.User())
			}
		}
		return ajaxSuccess(ctx, nil)
	}

	c.Prune = func(ctx web.Context) error {
		report, err := docker.VolumePrune()
		if err == nil {
			return ajaxSuccess(ctx, report)
		}
		return ajaxResult(ctx, err)
	}

	c.Detail = func(ctx web.Context) error {
		name := ctx.P("name")
		volume, _, err := docker.VolumeInspectRaw(name)
		if err != nil {
			return err
		}
		m := newModel(ctx).Add("Volume", volume)
		return ctx.Render("volume/detail", m)
	}

	c.Raw = func(ctx web.Context) error {
		name := ctx.P("name")
		_, raw, err := docker.VolumeInspectRaw(name)
		if err != nil {
			return err
		}

		j, err := misc.JSONIndent(raw)
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Volume", name).Add("Raw", j)
		return ctx.Render("volume/raw", m)
	}

	return
}

package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/model"
)

type ImageController struct {
	List web.HandlerFunc `path:"/" name:"image.list" authorize:"!" desc:"image list page"`
	//Detail web.HandlerFunc `path:"/:id/detail" name:"image.detail" authorize:"!" desc:"image detail page"`
	//Raw    web.HandlerFunc `path:"/:id/raw" name:"image.raw" authorize:"!" desc:"image raw page"`
}

func Image() (c *ImageController) {
	c = &ImageController{}

	c.List = func(ctx web.Context) error {
		name := ctx.Q("name")
		page := cast.ToIntD(ctx.Q("page"), 1)
		images, totalCount, err := docker.ImageList(name, page, model.PageSize)
		if err != nil {
			return err
		}

		m := newPagerModel(ctx, totalCount, model.PageSize, page).
			Add("Name", name).
			Add("Images", images)
		return ctx.Render("image/list", m)
	}

	return
}

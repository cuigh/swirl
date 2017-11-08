package controller

import (
	"strings"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

// ImageController is a controller of docker image
type ImageController struct {
	List   web.HandlerFunc `path:"/" name:"image.list" authorize:"!" desc:"image list page"`
	Detail web.HandlerFunc `path:"/:id/detail" name:"image.detail" authorize:"!" desc:"image detail page"`
	Raw    web.HandlerFunc `path:"/:id/raw" name:"image.raw" authorize:"!" desc:"image raw page"`
	Delete web.HandlerFunc `path:"/delete" method:"post" name:"image.delete" authorize:"!" desc:"delete image"`
}

// Image creates an instance of ImageController
func Image() (c *ImageController) {
	return &ImageController{
		List:   imageList,
		Detail: imageDetail,
		Raw:    imageRaw,
		Delete: imageDelete,
	}
}

func imageList(ctx web.Context) error {
	name := ctx.Q("name")
	page := cast.ToInt(ctx.Q("page"), 1)
	images, totalCount, err := docker.ImageList(name, page, model.PageSize)
	if err != nil {
		return err
	}

	m := newPagerModel(ctx, totalCount, model.PageSize, page).
		Set("Name", name).
		Set("Images", images)
	return ctx.Render("image/list", m)
}

func imageDetail(ctx web.Context) error {
	id := ctx.P("id")
	image, _, err := docker.ImageInspect(id)
	if err != nil {
		return err
	}

	histories, err := docker.ImageHistory(id)
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Image", image).Set("Histories", histories)
	return ctx.Render("image/detail", m)
}

func imageRaw(ctx web.Context) error {
	id := ctx.P("id")
	image, raw, err := docker.ImageInspect(id)
	if err != nil {
		return err
	}

	j, err := misc.JSONIndent(raw)
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Image", image).Set("Raw", j)
	return ctx.Render("image/raw", m)
}

func imageDelete(ctx web.Context) error {
	ids := strings.Split(ctx.F("ids"), ",")
	for _, id := range ids {
		if err := docker.ImageRemove(id); err != nil {
			return ajaxResult(ctx, err)
		}
	}
	return ajaxSuccess(ctx, nil)
}

package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/model"
)

func newModel(ctx web.Context) web.Map {
	return web.Map{
		"ContextUser": ctx.User(),
	}
}

func newPagerModel(ctx web.Context, totalCount, size, page int) web.Map {
	pager := model.NewPager(ctx.Request().RequestURI, totalCount, size, page)
	return newModel(ctx).Add("Pager", pager)
}

func ajaxResult(ctx web.Context, err error) error {
	if err != nil {
		return err
	}

	return ctx.JSON(web.Map{
		"success": err == nil,
	})
}

func ajaxSuccess(ctx web.Context, data interface{}) error {
	return ctx.JSON(web.Map{
		"success": true,
		"data":    data,
	})
}

package controller

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/model"
)

func newModel(ctx web.Context) data.Map {
	return data.Map{
		"ContextUser": ctx.User(),
	}
}

func newPagerModel(ctx web.Context, totalCount, size, page int) data.Map {
	pager := model.NewPager(ctx.Request().RequestURI, totalCount, size, page)
	return newModel(ctx).Set("Pager", pager)
}

func ajaxResult(ctx web.Context, err error) error {
	if err != nil {
		return err
	}

	return ctx.JSON(data.Map{
		"success": err == nil,
	})
}

func ajaxSuccess(ctx web.Context, value interface{}) error {
	return ctx.JSON(data.Map{
		"success": true,
		"data":    value,
	})
}

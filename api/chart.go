package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

// ChartHandler encapsulates chart related handlers.
type ChartHandler struct {
	Search web.HandlerFunc `path:"/search" auth:"chart.view" desc:"search charts"`
	Find   web.HandlerFunc `path:"/find" auth:"chart.view" desc:"find chart by id"`
	Save   web.HandlerFunc `path:"/save" method:"post" auth:"chart.edit" desc:"create or update chart"`
	Delete web.HandlerFunc `path:"/delete" method:"post" auth:"chart.delete" desc:"delete chart"`
}

// NewChart creates an instance of ChartHandler
func NewChart(b biz.ChartBiz) *ChartHandler {
	return &ChartHandler{
		Search: chartSearch(b),
		Find:   chartFind(b),
		Delete: chartDelete(b),
		Save:   chartSave(b),
	}
}

func chartSearch(b biz.ChartBiz) web.HandlerFunc {
	return func(ctx web.Context) (err error) {
		var (
			args   = &model.ChartSearchArgs{}
			charts []*model.Chart
			total  int
		)

		if err = ctx.Bind(args); err == nil {
			charts, total, err = b.Search(args)
		}

		if err != nil {
			return
		}

		return success(ctx, data.Map{
			"items": charts,
			"total": total,
		})
	}
}

func chartFind(b biz.ChartBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		id := ctx.Query("id")
		chart, err := b.Find(id)
		if err != nil {
			return err
		}
		return success(ctx, chart)
	}
}

func chartDelete(b biz.ChartBiz) web.HandlerFunc {
	type Args struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Delete(args.ID, args.Title, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func chartSave(b biz.ChartBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		r := &model.Chart{}
		err := ctx.Bind(r, true)
		if err == nil {
			if r.ID == "" {
				err = b.Create(r, ctx.User())
			} else {
				err = b.Update(r, ctx.User())
			}
		}
		return ajax(ctx, err)
	}
}

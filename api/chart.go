package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
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
	return func(c web.Context) (err error) {
		var (
			args   = &dao.ChartSearchArgs{}
			charts []*dao.Chart
			total  int
		)

		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			charts, total, err = b.Search(ctx, args)
		}

		if err != nil {
			return
		}

		return success(c, data.Map{
			"items": charts,
			"total": total,
		})
	}
}

func chartFind(b biz.ChartBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		id := c.Query("id")
		chart, err := b.Find(ctx, id)
		if err != nil {
			return err
		}
		return success(c, chart)
	}
}

func chartDelete(b biz.ChartBiz) web.HandlerFunc {
	type Args struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Delete(ctx, args.ID, args.Title, c.User())
		}
		return ajax(c, err)
	}
}

func chartSave(b biz.ChartBiz) web.HandlerFunc {
	return func(c web.Context) error {
		r := &dao.Chart{}
		err := c.Bind(r, true)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			if r.ID == "" {
				err = b.Create(ctx, r, c.User())
			} else {
				err = b.Update(ctx, r, c.User())
			}
		}
		return ajax(c, err)
	}
}

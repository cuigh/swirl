package api

import (
	"strings"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/ext/times"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

// ChartHandler encapsulates chart related handlers.
type ChartHandler struct {
	Search        web.HandlerFunc `path:"/search" auth:"chart.view" desc:"search charts"`
	Find          web.HandlerFunc `path:"/find" auth:"chart.view" desc:"find chart by id"`
	Save          web.HandlerFunc `path:"/save" method:"post" auth:"chart.edit" desc:"create or update chart"`
	Delete        web.HandlerFunc `path:"/delete" method:"post" auth:"chart.delete" desc:"delete chart"`
	FetchData     web.HandlerFunc `path:"/fetch-data" auth:"?" desc:"fetch chart data"`
	FindDashboard web.HandlerFunc `path:"/find-dashboard" auth:"?" desc:"find dashboard by name and key"`
	SaveDashboard web.HandlerFunc `path:"/save-dashboard" method:"post" auth:"chart.dashboard" desc:"save dashboard"`
}

// NewChart creates an instance of ChartHandler
func NewChart(b biz.ChartBiz) *ChartHandler {
	return &ChartHandler{
		Search:        chartSearch(b),
		Find:          chartFind(b),
		Delete:        chartDelete(b),
		Save:          chartSave(b),
		FetchData:     chartFetchData(b),
		FindDashboard: chartFindDashboard(b),
		SaveDashboard: chartSaveDashboard(b),
	}
}

func chartSearch(b biz.ChartBiz) web.HandlerFunc {
	return func(ctx web.Context) (err error) {
		var (
			args   = &model.ChartSearchArgs{}
			charts []*biz.Chart
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
		r := &biz.Chart{}
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

func chartFetchData(b biz.ChartBiz) web.HandlerFunc {
	type Args struct {
		Key    string `json:"key" bind:"key"`
		Charts string `json:"charts" bind:"charts"`
		Period int32  `json:"period" bind:"period"`
	}
	return func(ctx web.Context) (err error) {
		var (
			args = &Args{}
			d    data.Map
		)
		if err = ctx.Bind(args); err == nil {
			d, err = b.FetchData(args.Key, strings.Split(args.Charts, ","), times.Minutes(args.Period))
		}
		if err != nil {
			return err
		}
		return success(ctx, d)
	}
}

func chartFindDashboard(b biz.ChartBiz) web.HandlerFunc {
	return func(ctx web.Context) (err error) {
		var (
			d    *biz.Dashboard
			name = ctx.Query("name")
			key  = ctx.Query("key")
		)
		d, err = b.FindDashboard(name, key)
		if err != nil {
			return err
		}
		return success(ctx, d)
	}
}

func chartSaveDashboard(b biz.ChartBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		dashboard := &model.Dashboard{}
		err := ctx.Bind(dashboard)
		if err != nil {
			return err
		}

		switch dashboard.Name {
		case "home", "service":
			err = b.UpdateDashboard(dashboard, ctx.User())
		default:
			err = errors.New("unknown dashboard: " + dashboard.Name)
		}
		return ajax(ctx, err)
	}
}

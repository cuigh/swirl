package api

import (
	"strings"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/ext/times"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/dao"
)

// DashboardHandler encapsulates dashboard related handlers.
type DashboardHandler struct {
	Find      web.HandlerFunc `path:"/find" auth:"?" desc:"find dashboard by name and key"`
	Save      web.HandlerFunc `path:"/save" method:"post" auth:"dashboard.edit" desc:"save dashboard"`
	FetchData web.HandlerFunc `path:"/fetch-data" auth:"?" desc:"fetch data of dashboard charts"`
}

// NewDashboard creates an instance of DashboardHandler
func NewDashboard(b biz.DashboardBiz) *DashboardHandler {
	return &DashboardHandler{
		Find:      dashboardFind(b),
		Save:      dashboardSave(b),
		FetchData: dashboardFetchData(b),
	}
}

func dashboardFind(b biz.DashboardBiz) web.HandlerFunc {
	return func(ctx web.Context) (err error) {
		var (
			d    *dao.Dashboard
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

func dashboardSave(b biz.DashboardBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		dashboard := &dao.Dashboard{}
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

func dashboardFetchData(b biz.DashboardBiz) web.HandlerFunc {
	type Args struct {
		Key        string `json:"key" bind:"key"`
		Dashboards string `json:"charts" bind:"charts"`
		Period     int32  `json:"period" bind:"period"`
	}
	return func(ctx web.Context) (err error) {
		var (
			args = &Args{}
			d    data.Map
		)
		if err = ctx.Bind(args); err == nil {
			d, err = b.FetchData(args.Key, strings.Split(args.Dashboards, ","), times.Minutes(args.Period))
		}
		if err != nil {
			return err
		}
		return success(ctx, d)
	}
}

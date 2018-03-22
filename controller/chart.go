package controller

import (
	"time"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

// ChartController is a controller of metric chart.
type ChartController struct {
	List   web.HandlerFunc `path:"/" name:"chart.list" authorize:"!" desc:"chart list page"`
	Query  web.HandlerFunc `path:"/query" name:"chart.query" authorize:"?" desc:"chart query"`
	New    web.HandlerFunc `path:"/new" name:"chart.new" authorize:"!" desc:"new chart page"`
	Create web.HandlerFunc `path:"/new" method:"post" name:"chart.create" authorize:"!" desc:"create chart"`
	Edit   web.HandlerFunc `path:"/:name/edit" name:"chart.edit" authorize:"!" desc:"edit chart page"`
	Delete web.HandlerFunc `path:"/:name/delete" method:"post" name:"chart.delete" authorize:"!" desc:"delete chart"`
	Update web.HandlerFunc `path:"/:name/edit" method:"post" name:"chart.update" authorize:"!" desc:"update chart"`
	Data   web.HandlerFunc `path:"/data" name:"chart.data" authorize:"?" desc:"fetch chart datas"`
}

// Chart creates an instance of RoleController
func Chart() (c *ChartController) {
	return &ChartController{
		List:   chartList,
		Query:  chartQuery,
		New:    chartNew,
		Create: chartCreate,
		Edit:   chartEdit,
		Update: chartUpdate,
		Delete: chartDelete,
		Data:   chartData,
	}
}

func chartList(ctx web.Context) error {
	charts, err := biz.Chart.List()
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Charts", charts)
	return ctx.Render("system/chart/list", m)
}

func chartQuery(ctx web.Context) error {
	charts, err := biz.Chart.List()
	if err != nil {
		return err
	}

	dashboard := ctx.Q("dashboard")
	var list []*model.Chart
	for _, c := range charts {
		if c.Dashboard == dashboard || c.Dashboard == "" {
			list = append(list, c)
		}
	}
	return ctx.JSON(list)
}

func chartNew(ctx web.Context) error {
	m := newModel(ctx).Set("Chart", &model.Chart{
		Width:     12,
		Height:    50,
		Type:      "line",
		Dashboard: "service",
	})
	return ctx.Render("system/chart/edit", m)
}

func chartCreate(ctx web.Context) error {
	chart := &model.Chart{}
	err := ctx.Bind(chart, true)
	if err == nil {
		err = biz.Chart.Create(chart, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func chartEdit(ctx web.Context) error {
	name := ctx.P("name")
	chart, err := biz.Chart.Get(name)
	if err != nil {
		return err
	}
	if chart == nil {
		return web.ErrNotFound
	}

	m := newModel(ctx).Set("Chart", chart)
	return ctx.Render("system/chart/edit", m)
}

func chartUpdate(ctx web.Context) error {
	chart := &model.Chart{}
	err := ctx.Bind(chart)
	if err == nil {
		err = biz.Chart.Update(chart, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func chartDelete(ctx web.Context) error {
	name := ctx.P("name")
	err := biz.Chart.Delete(name, ctx.User())
	return ajaxResult(ctx, err)
}

// todo:
func chartData(ctx web.Context) error {
	period := cast.ToDuration(ctx.Q("time"), time.Hour)
	dashboard := ctx.Q("dashboard")

	var (
		charts []*model.Chart
		err    error
	)
	if dashboard == "home" {
		var setting *model.Setting
		setting, err = biz.Setting.Get()
		if err != nil {
			return err
		}
		charts, err = biz.Chart.Panel(setting.Dashboard.Home)

	} else if dashboard == "service" {
		id := ctx.Q("id")
		charts, err = biz.Chart.GetServiceCharts(id)
	}

	if err != nil {
		return err
	}

	datas, err := biz.Chart.FetchDatas(charts, period)
	if err != nil {
		return err
	}
	return ctx.JSON(datas)
}

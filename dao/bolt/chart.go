package bolt

import (
	"context"
	"sort"

	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

func (d *Dao) ChartList(ctx context.Context, title, dashboard string, pageIndex, pageSize int) (charts []*model.Chart, count int, err error) {
	err = d.each("chart", func(v Value) error {
		chart := &model.Chart{}
		err = v.Unmarshal(chart)
		if err == nil {
			match := true
			if title != "" {
				match = matchAny(title, chart.Title)
			}
			if match && dashboard != "" {
				match = matchAny(chart.Dashboard, dashboard, "")
			}
			if match {
				charts = append(charts, chart)
			}
		}
		return err
	})
	if err == nil {
		count = len(charts)
		sort.Slice(charts, func(i, j int) bool {
			return charts[i].CreatedAt.After(charts[j].UpdatedAt)
		})
		start, end := misc.Page(count, pageIndex, pageSize)
		charts = charts[start:end]
	}
	return
}

func (d *Dao) ChartCreate(ctx context.Context, chart *model.Chart) (err error) {
	return d.update("chart", chart.ID, chart)
}

func (d *Dao) ChartGet(ctx context.Context, name string) (chart *model.Chart, err error) {
	var v Value
	v, err = d.get("chart", name)
	if err == nil && v != nil {
		chart = &model.Chart{}
		err = v.Unmarshal(chart)
	}
	return
}

func (d *Dao) ChartBatch(ctx context.Context, names ...string) (charts []*model.Chart, err error) {
	err = d.slice("chart", func(v Value) error {
		chart := &model.Chart{}
		err = v.Unmarshal(chart)
		if err == nil {
			charts = append(charts, chart)
		}
		return err
	}, names...)
	return
}

func (d *Dao) ChartUpdate(ctx context.Context, chart *model.Chart) (err error) {
	return d.update("chart", chart.ID, chart)
}

func (d *Dao) ChartDelete(ctx context.Context, name string) (err error) {
	return d.delete("chart", name)
}

func (d *Dao) DashboardGet(ctx context.Context, name, key string) (dashboard *model.Dashboard, err error) {
	cd := &model.Dashboard{
		Name: name,
		Key:  key,
	}

	var v Value
	v, err = d.get("dashboard", cd.ID())
	if v != nil {
		if err = v.Unmarshal(cd); err == nil {
			return cd, nil
		}
	}
	return nil, err
}

func (d *Dao) DashboardUpdate(ctx context.Context, dashboard *model.Dashboard) (err error) {
	return d.update("dashboard", dashboard.ID(), dashboard)
}

package bolt

import (
	"context"
	"sort"
	"time"

	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
)

const (
	Chart     = "chart"
	Dashboard = "dashboard"
)

func (d *Dao) ChartSearch(ctx context.Context, args *dao.ChartSearchArgs) (charts []*dao.Chart, count int, err error) {
	err = d.each(Chart, func(v []byte) error {
		chart := &dao.Chart{}
		if err = decode(v, chart); err == nil {
			match := true
			if args.Title != "" {
				match = matchAny(args.Title, chart.Title)
			}
			if match && args.Dashboard != "" {
				match = matchAny(chart.Dashboard, args.Dashboard, "")
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
			return time.Time(charts[i].UpdatedAt).After(time.Time(charts[j].UpdatedAt))
		})
		start, end := misc.Page(count, args.PageIndex, args.PageSize)
		charts = charts[start:end]
	}
	return
}

func (d *Dao) ChartCreate(ctx context.Context, chart *dao.Chart) (err error) {
	return d.replace(Chart, chart.ID, chart)
}

func (d *Dao) ChartGet(ctx context.Context, name string) (chart *dao.Chart, err error) {
	chart = &dao.Chart{}
	err = d.get(Chart, name, chart)
	if err == ErrNoRecords {
		return nil, nil
	} else if err != nil {
		chart = nil
	}
	return
}

func (d *Dao) ChartGetBatch(ctx context.Context, ids ...string) (charts []*dao.Chart, err error) {
	err = d.slice(Chart, func(v []byte) error {
		chart := &dao.Chart{}
		if err = decode(v, chart); err == nil {
			charts = append(charts, chart)
		}
		return err
	}, ids...)
	return
}

func (d *Dao) ChartUpdate(ctx context.Context, chart *dao.Chart) (err error) {
	old := &dao.Chart{}
	return d.update(Chart, chart.ID, old, func() interface{} {
		chart.CreatedAt = old.CreatedAt
		chart.CreatedBy = old.CreatedBy
		return chart
	})
}

func (d *Dao) ChartDelete(ctx context.Context, name string) (err error) {
	return d.delete(Chart, name)
}

func (d *Dao) DashboardGet(ctx context.Context, name, key string) (dashboard *dao.Dashboard, err error) {
	dashboard = &dao.Dashboard{
		Name: name,
		Key:  key,
	}
	err = d.get(Dashboard, dashboard.ID(), dashboard)
	if err == ErrNoRecords {
		return nil, nil
	} else if err != nil {
		dashboard = nil
	}
	return
}

func (d *Dao) DashboardUpdate(ctx context.Context, dashboard *dao.Dashboard) (err error) {
	return d.replace(Dashboard, dashboard.ID(), dashboard)
}

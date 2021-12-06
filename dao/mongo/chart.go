package mongo

import (
	"context"

	"github.com/cuigh/swirl/model"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	Chart     = "chart"
	Dashboard = "dashboard"
)

func (d *Dao) ChartList(ctx context.Context, title, dashboard string, pageIndex, pageSize int) (charts []*model.Chart, count int, err error) {
	filter := bson.M{}
	if title != "" {
		filter["title"] = title
	}
	if dashboard != "" {
		filter["dashboard"] = bson.M{"$in": []string{"", dashboard}}
	}
	opts := searchOptions{filter: filter, sorter: bson.M{"updated_at": -1}, pageIndex: pageIndex, pageSize: pageSize}
	charts = []*model.Chart{}
	count, err = d.search(ctx, Chart, opts, &charts)
	return
}

func (d *Dao) ChartCreate(ctx context.Context, chart *model.Chart) (err error) {
	return d.create(ctx, Chart, chart)
}

func (d *Dao) ChartGet(ctx context.Context, id string) (chart *model.Chart, err error) {
	chart = &model.Chart{}
	found, err := d.find(ctx, Chart, id, chart)
	if !found {
		return nil, err
	}
	return
}

func (d *Dao) ChartBatch(ctx context.Context, names ...string) (charts []*model.Chart, err error) {
	charts = []*model.Chart{}
	err = d.fetch(ctx, Chart, bson.M{"_id": bson.M{"$in": names}}, &charts)
	return
}

func (d *Dao) ChartUpdate(ctx context.Context, chart *model.Chart) (err error) {
	return d.update(ctx, Chart, chart.ID, chart)
}

func (d *Dao) ChartDelete(ctx context.Context, id string) (err error) {
	return d.delete(ctx, Chart, id)
}

func (d *Dao) DashboardGet(ctx context.Context, name, key string) (dashboard *model.Dashboard, err error) {
	dashboard = &model.Dashboard{
		Name: name,
		Key:  key,
	}
	found, err := d.find(ctx, Dashboard, dashboard.ID(), dashboard)
	if !found {
		return nil, err
	}
	return
}

func (d *Dao) DashboardUpdate(ctx context.Context, dashboard *model.Dashboard) (err error) {
	update := bson.M{
		"$set": dashboard,
	}
	return d.update(ctx, Dashboard, dashboard.ID(), update)
}

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

func (d *Dao) ChartSearch(ctx context.Context, args *model.ChartSearchArgs) (charts []*model.Chart, count int, err error) {
	filter := bson.M{}
	if args.Title != "" {
		filter["title"] = args.Title
	}
	if args.Dashboard != "" {
		filter["dashboard"] = bson.M{"$in": []string{"", args.Dashboard}}
	}
	opts := searchOptions{filter: filter, sorter: bson.M{"updated_at": -1}, pageIndex: args.PageIndex, pageSize: args.PageSize}
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

func (d *Dao) ChartGetBatch(ctx context.Context, names ...string) (charts []*model.Chart, err error) {
	charts = []*model.Chart{}
	err = d.fetch(ctx, Chart, bson.M{"_id": bson.M{"$in": names}}, &charts)
	return
}

func (d *Dao) ChartUpdate(ctx context.Context, chart *model.Chart) (err error) {
	update := bson.M{
		"$set": bson.M{
			"title":      chart.Title,
			"desc":       chart.Description,
			"width":      chart.Width,
			"height":     chart.Height,
			"unit":       chart.Unit,
			"dashboard":  chart.Dashboard,
			"type":       chart.Type,
			"margin":     chart.Margin,
			"metrics":    chart.Metrics,
			"updated_at": chart.UpdatedAt,
			"updated_by": chart.UpdatedBy,
		},
	}
	return d.update(ctx, Chart, chart.ID, update)
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
	return d.upsert(ctx, Dashboard, dashboard.ID(), update)
}

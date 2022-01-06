package biz

import (
	"context"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
)

type ChartBiz interface {
	Search(ctx context.Context, args *dao.ChartSearchArgs) (charts []*dao.Chart, total int, err error)
	Delete(ctx context.Context, id, title string, user web.User) (err error)
	Find(ctx context.Context, id string) (chart *dao.Chart, err error)
	Create(ctx context.Context, chart *dao.Chart, user web.User) (err error)
	Update(ctx context.Context, chart *dao.Chart, user web.User) (err error)
}

func NewChart(d dao.Interface, mb MetricBiz, eb EventBiz) ChartBiz {
	return &chartBiz{
		d:  d,
		mb: mb,
		eb: eb,
	}
}

type chartBiz struct {
	d  dao.Interface
	mb MetricBiz
	eb EventBiz
}

func (b *chartBiz) Search(ctx context.Context, args *dao.ChartSearchArgs) (charts []*dao.Chart, total int, err error) {
	return b.d.ChartSearch(ctx, args)
}

func (b *chartBiz) Create(ctx context.Context, chart *dao.Chart, user web.User) (err error) {
	chart.ID = createId()
	chart.CreatedAt = now()
	chart.CreatedBy = newOperator(user)
	chart.UpdatedAt = chart.CreatedAt
	chart.UpdatedBy = chart.CreatedBy
	err = b.d.ChartCreate(ctx, chart)
	if err == nil {
		b.eb.CreateChart(EventActionCreate, chart.ID, chart.Title, user)
	}
	return
}

func (b *chartBiz) Delete(ctx context.Context, id, title string, user web.User) (err error) {
	err = b.d.ChartDelete(ctx, id)
	if err == nil {
		b.eb.CreateChart(EventActionDelete, id, title, user)
	}
	return
}

func (b *chartBiz) Find(ctx context.Context, id string) (chart *dao.Chart, err error) {
	return b.d.ChartGet(ctx, id)
}

func (b *chartBiz) Update(ctx context.Context, chart *dao.Chart, user web.User) (err error) {
	chart.UpdatedAt = now()
	chart.UpdatedBy = newOperator(user)
	err = b.d.ChartUpdate(ctx, chart)
	if err == nil {
		b.eb.CreateChart(EventActionUpdate, chart.ID, chart.Title, user)
	}
	return
}

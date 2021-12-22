package biz

import (
	"context"
	"os"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
	"github.com/jinzhu/copier"
)

var builtins = []*model.Chart{
	model.NewChart("service", "$cpu", "CPU", "${name}", `rate(container_cpu_user_seconds_total{container_label_com_docker_swarm_service_name="${service}"}[5m]) * 100`, "percent:100", 60),
	model.NewChart("service", "$memory", "Memory", "${name}", `container_memory_usage_bytes{container_label_com_docker_swarm_service_name="${service}"}`, "size:bytes", 60),
	model.NewChart("service", "$network_in", "Network Receive", "${name}", `sum(irate(container_network_receive_bytes_total{container_label_com_docker_swarm_service_name="${service}"}[5m])) by(name)`, "size:bytes", 60),
	model.NewChart("service", "$network_out", "Network Send", "${name}", `sum(irate(container_network_transmit_bytes_total{container_label_com_docker_swarm_service_name="${service}"}[5m])) by(name)`, "size:bytes", 60),
}

type DashboardBiz interface {
	FetchData(key string, ids []string, period time.Duration) (data.Map, error)
	FindDashboard(name, key string) (dashboard *model.Dashboard, err error)
	UpdateDashboard(dashboard *model.Dashboard, user web.User) (err error)
}

func NewDashboard(d dao.Interface, mb MetricBiz, eb EventBiz) DashboardBiz {
	return &dashboardBiz{
		d:  d,
		mb: mb,
		eb: eb,
	}
}

type dashboardBiz struct {
	d  dao.Interface
	cb ChartBiz
	mb MetricBiz
	eb EventBiz
}

func (b *dashboardBiz) FindDashboard(name, key string) (dashboard *model.Dashboard, err error) {
	if dashboard, err = b.d.DashboardGet(context.TODO(), name, key); err != nil {
		return
	}
	if dashboard == nil {
		dashboard = b.defaultDashboard(name, key)
	}
	err = b.fillCharts(dashboard)
	return
}

func (b *dashboardBiz) UpdateDashboard(dashboard *model.Dashboard, user web.User) (err error) {
	dashboard.UpdatedAt = now()
	dashboard.UpdatedBy = newOperator(user)
	return b.d.DashboardUpdate(context.TODO(), dashboard)
}

func (b *dashboardBiz) FetchData(key string, ids []string, period time.Duration) (data.Map, error) {
	if !b.mb.Enabled() {
		return data.Map{}, nil
	}

	charts, err := b.getCharts(ids)
	if err != nil {
		return nil, err
	}

	type Data struct {
		id   string
		data interface{}
		err  error
	}

	ch := make(chan Data, len(charts))
	end := time.Now()
	start := end.Add(-period)
	for _, chart := range charts {
		go func(c *model.Chart) {
			d := Data{id: c.ID}
			switch c.Type {
			case "line", "bar":
				d.data, d.err = b.fetchMatrixData(c, key, start, end)
			case "pie":
				d.data, d.err = b.fetchVectorData(c, key, end)
			case "gauge":
				d.data, d.err = b.fetchScalarData(c, key, end)
			default:
				d.err = errors.New("invalid chart type: " + c.Type)
			}
			ch <- d
		}(chart)
	}

	ds := data.Map{}
	for range charts {
		d := <-ch
		if d.err != nil {
			log.Get("metric").Error(d.err)
		} else {
			ds.Set(d.id, d.data)
		}
	}
	close(ch)
	return ds, nil
}

func (b *dashboardBiz) fetchMatrixData(chart *model.Chart, key string, start, end time.Time) (md *MatrixData, err error) {
	var (
		q string
		d *MatrixData
	)
	for i, m := range chart.Metrics {
		q, err = b.formatQuery(m.Query, chart.Dashboard, key)
		if err != nil {
			return nil, err
		}

		if d, err = b.mb.GetMatrix(q, m.Legend, start, end); err != nil {
			log.Get("metric").Error(err)
		} else if i == 0 {
			md = d
		} else {
			md.Legend = append(md.Legend, d.Legend...)
			md.Series = append(md.Series, d.Series...)
		}
	}
	return md, nil
}

func (b *dashboardBiz) fetchVectorData(chart *model.Chart, key string, end time.Time) (cvd *VectorData, err error) {
	var (
		q string
		d *VectorData
	)
	for i, m := range chart.Metrics {
		q, err = b.formatQuery(m.Query, chart.Dashboard, key)
		if err != nil {
			return nil, err
		}

		if d, err = b.mb.GetVector(q, m.Legend, end); err != nil {
			log.Get("metric").Error(err)
		} else if i == 0 {
			cvd = d
		} else {
			cvd.Legend = append(cvd.Legend, d.Legend...)
			cvd.Data = append(cvd.Data, d.Data...)
		}
	}
	return cvd, nil
}

func (b *dashboardBiz) fetchScalarData(chart *model.Chart, key string, end time.Time) (*VectorValue, error) {
	query, err := b.formatQuery(chart.Metrics[0].Query, chart.Dashboard, key)
	if err != nil {
		return nil, err
	}

	v, err := b.mb.GetScalar(query, end)
	if err != nil {
		return nil, err
	}

	return &VectorValue{
		//Name:  "",
		Value: v,
	}, nil
}

func (b *dashboardBiz) formatQuery(query, dashboard, key string) (string, error) {
	if dashboard == "home" {
		return query, nil
	}

	var errs []error
	m := map[string]string{dashboard: key}
	q := os.Expand(query, func(k string) string {
		if v, ok := m[k]; ok {
			return v
		}
		errs = append(errs, errors.New("invalid argument in query: "+query))
		return ""
	})
	if len(errs) == 0 {
		return q, nil
	}
	return "", errs[0]
}

func (b *dashboardBiz) getCharts(ids []string) (charts map[string]*model.Chart, err error) {
	var (
		customIds    []string
		customCharts []*model.Chart
	)

	charts = make(map[string]*model.Chart)
	for _, id := range ids {
		if id[0] == '$' {
			for _, c := range builtins {
				if c.ID == id {
					charts[id] = c
				}
			}
		} else {
			customIds = append(customIds, id)
		}
	}

	if len(customIds) > 0 {
		if customCharts, err = b.d.ChartGetBatch(context.TODO(), customIds...); err == nil {
			for _, chart := range customCharts {
				charts[chart.ID] = chart
			}
		}
	}
	return
}

func (b *dashboardBiz) fillCharts(d *model.Dashboard) (err error) {
	if len(d.Charts) == 0 {
		return
	}

	var (
		m   map[string]*model.Chart
		ids = make([]string, len(d.Charts))
	)

	for i, c := range d.Charts {
		ids[i] = c.ID
	}

	m, err = b.getCharts(ids)
	if err != nil {
		return err
	}

	for i := range d.Charts {
		if c := m[d.Charts[i].ID]; c != nil {
			_ = copier.CopyWithOption(&d.Charts[i], c, copier.Option{IgnoreEmpty: true})
		}
	}
	return nil
}

func (b *dashboardBiz) defaultDashboard(name, key string) *model.Dashboard {
	d := &model.Dashboard{
		Name:     name,
		Key:      key,
		Period:   30,
		Interval: 15,
	}
	if name == "service" {
		d.Charts = []model.ChartInfo{
			{ID: "$cpu"},
			{ID: "$memory"},
			{ID: "$network_in"},
			{ID: "$network_out"},
		}
	}
	return d
}

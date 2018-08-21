package biz

import (
	"os"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
)

// Chart return a chart biz instance.
var Chart = newChartBiz()

type chartBiz struct {
	builtin []*model.Chart
}

func newChartBiz() *chartBiz {
	b := &chartBiz{}
	b.builtin = append(b.builtin, model.NewChart("service", "$cpu", "CPU", "${name}", `rate(container_cpu_user_seconds_total{container_label_com_docker_swarm_service_name="${service}"}[5m]) * 100`, "percent:100"))
	b.builtin = append(b.builtin, model.NewChart("service", "$memory", "Memory", "${name}", `container_memory_usage_bytes{container_label_com_docker_swarm_service_name="${service}"}`, "size:bytes"))
	b.builtin = append(b.builtin, model.NewChart("service", "$network_in", "Network Receive", "${name}", `sum(irate(container_network_receive_bytes_total{container_label_com_docker_swarm_service_name="${service}"}[5m])) by(name)`, "size:bytes"))
	b.builtin = append(b.builtin, model.NewChart("service", "$network_out", "Network Send", "${name}", `sum(irate(container_network_transmit_bytes_total{container_label_com_docker_swarm_service_name="${service}"}[5m])) by(name)`, "size:bytes"))
	return b
}

func (b *chartBiz) List() (charts []*model.Chart, err error) {
	do(func(d dao.Interface) {
		charts, err = d.ChartList()
	})
	return
}

func (b *chartBiz) Create(chart *model.Chart, user web.User) (err error) {
	do(func(d dao.Interface) {
		// chart.CreatedAt = time.Now()
		// chart.UpdatedAt = chart.CreatedAt
		err = d.ChartCreate(chart)
	})
	return
}

func (b *chartBiz) Delete(id string, user web.User) (err error) {
	do(func(d dao.Interface) {
		err = d.ChartDelete(id)
	})
	return
}

func (b *chartBiz) Get(name string) (chart *model.Chart, err error) {
	do(func(d dao.Interface) {
		chart, err = d.ChartGet(name)
		if len(chart.Metrics) == 0 {
			chart.Metrics = append(chart.Metrics, model.ChartMetric{Legend: chart.Legend, Query: chart.Query})
		}
	})
	return
}

func (b *chartBiz) Batch(names ...string) (charts []*model.Chart, err error) {
	do(func(d dao.Interface) {
		charts, err = d.ChartBatch(names...)
		if err == nil {
			for _, c := range charts {
				if len(c.Metrics) == 0 {
					c.Metrics = append(c.Metrics, model.ChartMetric{Legend: c.Legend, Query: c.Query})
				}
			}
		}
	})
	return
}

func (b *chartBiz) Update(chart *model.Chart, user web.User) (err error) {
	do(func(d dao.Interface) {
		// chart.UpdatedAt = time.Now()
		err = d.ChartUpdate(chart)
	})
	return
}

func (b *chartBiz) GetServiceCharts(name string) (charts []*model.Chart, err error) {
	// service, _, err := docker.ServiceInspect(name)
	// if err != nil {
	// 	return nil, err
	// }

	// if label := service.Spec.Labels["swirl.metrics"]; label != "" {
	// 	names := strings.Split(label, ",")
	// }
	charts = b.builtin
	return
}

func (b *chartBiz) GetDashboard(name, key string) (dashboard *model.ChartDashboard, err error) {
	do(func(d dao.Interface) {
		dashboard, err = d.DashboardGet(name, key)
	})
	return
}

func (b *chartBiz) UpdateDashboard(dashboard *model.ChartDashboard, user web.User) (err error) {
	do(func(d dao.Interface) {
		err = d.DashboardUpdate(dashboard)
	})
	return
}

// nolint: gocyclo
func (b *chartBiz) GetDashboardCharts(dashboard *model.ChartDashboard) (charts []*model.Chart, err error) {
	do(func(d dao.Interface) {
		if dashboard == nil || len(dashboard.Charts) == 0 {
			return
		}

		names := make([]string, len(dashboard.Charts))
		for i, c := range dashboard.Charts {
			names[i] = c.Name
		}

		var cs []*model.Chart
		cs, err = b.getCharts(names)
		if err != nil {
			return
		}

		if len(cs) > 0 {
			m := make(map[string]*model.Chart)
			for _, c := range cs {
				m[c.Name] = c
			}
			for _, c := range dashboard.Charts {
				if chart := m[c.Name]; chart != nil {
					if c.Width > 0 {
						chart.Width = c.Width
					}
					if c.Height > 0 {
						chart.Height = c.Height
					}
					//if len(c.Colors) > 0 {
					//	chart.Colors = c.Colors
					//}
					charts = append(charts, chart)
				}
			}
		}
	})
	return
}

func (b *chartBiz) FetchDatas(key string, names []string, period time.Duration) (data.Map, error) {
	charts, err := b.getCharts(names)
	if err != nil {
		return nil, err
	}

	type Data struct {
		name string
		data interface{}
		err  error
	}

	ch := make(chan Data, len(charts))
	end := time.Now()
	start := end.Add(-period)
	for _, chart := range charts {
		go func(c *model.Chart) {
			d := Data{name: c.Name}
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
			ds.Set(d.name, d.data)
		}
	}
	close(ch)
	return ds, nil
}

func (b *chartBiz) fetchMatrixData(chart *model.Chart, key string, start, end time.Time) (*model.ChartMatrixData, error) {
	var cmd *model.ChartMatrixData
	for i, m := range chart.Metrics {
		q, err := b.formatQuery(m.Query, chart.Dashboard, key)
		if err != nil {
			return nil, err
		}

		if d, err := Metric.GetMatrix(q, m.Legend, start, end); err != nil {
			log.Get("metric").Error(err)
		} else if i == 0 {
			cmd = d
		} else {
			cmd.Legend = append(cmd.Legend, d.Legend...)
			cmd.Series = append(cmd.Series, d.Series...)
		}
	}
	return cmd, nil
}

func (b *chartBiz) fetchVectorData(chart *model.Chart, key string, end time.Time) (*model.ChartVectorData, error) {
	var cvd *model.ChartVectorData
	for i, m := range chart.Metrics {
		query, err := b.formatQuery(m.Query, chart.Dashboard, key)
		if err != nil {
			return nil, err
		}

		if d, err := Metric.GetVector(query, m.Legend, end); err != nil {
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

func (b *chartBiz) fetchScalarData(chart *model.Chart, key string, end time.Time) (*model.ChartValue, error) {
	query, err := b.formatQuery(chart.Metrics[0].Query, chart.Dashboard, key)
	if err != nil {
		return nil, err
	}

	v, err := Metric.GetScalar(query, end)
	if err != nil {
		return nil, err
	}

	return &model.ChartValue{
		//Name:  "",
		Value: v,
	}, nil
}

func (b *chartBiz) formatQuery(query, dashboard, key string) (string, error) {
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

func (b *chartBiz) getCharts(names []string) (charts []*model.Chart, err error) {
	var (
		customNames  []string
		customCharts []*model.Chart
	)

	for _, n := range names {
		if n[0] == '$' {
			for _, c := range b.builtin {
				if c.Name == n {
					charts = append(charts, c)
				}
			}
		} else {
			customNames = append(customNames, n)
		}
	}

	if len(customNames) > 0 {
		if customCharts, err = b.Batch(customNames...); err == nil {
			charts = append(charts, customCharts...)
		}
	}
	return
}

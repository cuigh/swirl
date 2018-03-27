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
	// for _, c := range categories {
	// 	if c == "java" {
	// 		charts = append(charts, model.NewChart("threads", "Threads", "${instance}", `jvm_threads_current{service="%s"}`, ""))
	// 		charts = append(charts, model.NewChart("gc_duration", "GC Duration", "${instance}", `rate(jvm_gc_collection_seconds_sum{service="%s"}[1m])`, "time:s"))
	// 	} else if c == "go" {
	// 		charts = append(charts, model.NewChart("threads", "Threads", "${instance}", `go_threads{service="%s"}`, ""))
	// 		charts = append(charts, model.NewChart("goroutines", "Goroutines", "${instance}", `go_goroutines{service="%s"}`, ""))
	// 		charts = append(charts, model.NewChart("gc_duration", "GC Duration", "${instance}", `sum(go_gc_duration_seconds{service="%s"}) by (instance)`, "time:s"))
	// 	}
	// }
	// for i, c := range charts {
	// 	charts[i].Query = fmt.Sprintf(c.Query, name)
	// }
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
	})
	return
}

func (b *chartBiz) Batch(names ...string) (charts []*model.Chart, err error) {
	do(func(d dao.Interface) {
		charts, err = d.ChartBatch(names...)
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
		if len(dashboard.Charts) == 0 {
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
					if len(c.Colors) > 0 {
						chart.Colors = c.Colors
					}
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

	datas := data.Map{}
	end := time.Now()
	start := end.Add(-period)
	for _, chart := range charts {
		query, err := b.formatQuery(chart, key)
		if err != nil {
			return nil, err
		}

		switch chart.Type {
		case "line", "bar":
			m, err := Metric.GetMatrix(query, chart.Label, start, end)
			if err != nil {
				log.Get("metric").Error(err, query)
			} else {
				datas[chart.Name] = m
			}
		case "pie", "table":
			m, err := Metric.GetVector(query, chart.Label, end)
			if err != nil {
				log.Get("metric").Error(err, query)
			} else {
				datas[chart.Name] = m
			}
		case "value":
		}
	}
	return datas, nil
}

func (b *chartBiz) formatQuery(chart *model.Chart, key string) (string, error) {
	if chart.Dashboard == "home" {
		return chart.Query, nil
	}

	var errs []error
	m := map[string]string{chart.Dashboard: key}
	query := os.Expand(chart.Query, func(k string) string {
		if v, ok := m[k]; ok {
			return v
		}
		errs = append(errs, errors.New("invalid argument in query: "+chart.Query))
		return ""
	})
	if len(errs) == 0 {
		return query, nil
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

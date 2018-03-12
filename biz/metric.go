package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/cuigh/auxo/ext/times"
	"github.com/cuigh/auxo/util/lazy"
	"github.com/cuigh/swirl/model"
	pclient "github.com/prometheus/client_golang/api"
	papi "github.com/prometheus/client_golang/api/prometheus/v1"
	pmodel "github.com/prometheus/common/model"
)

// Metric return a metric biz instance.
var Metric = &metricBiz{
	api: lazy.Value{
		New: func() (interface{}, error) {
			setting, err := Setting.Get()
			if err != nil {
				return nil, err
			}

			client, err := pclient.NewClient(pclient.Config{Address: setting.Metrics.Prometheus})
			if err != nil {
				return nil, err
			}

			return papi.NewAPI(client), nil
		},
	},
}

type metricBiz struct {
	api lazy.Value
}

func (b *metricBiz) GetServiceCharts(service string, categories []string) (charts []model.ChartInfo) {
	charts = append(charts, model.NewChartInfo("cpu", "CPU", "name", `rate(container_cpu_user_seconds_total{container_label_com_docker_swarm_service_name="%s"}[5m]) * 100`))
	charts = append(charts, model.NewChartInfo("memory", "Memory", "name", `container_memory_usage_bytes{container_label_com_docker_swarm_service_name="%s"}`))
	charts = append(charts, model.NewChartInfo("network_in", "Network Receive", "name", `sum(irate(container_network_receive_bytes_total{container_label_com_docker_swarm_service_name="%s"}[5m])) by(name)`))
	charts = append(charts, model.NewChartInfo("network_out", "Network Send", "name", `sum(irate(container_network_transmit_bytes_total{container_label_com_docker_swarm_service_name="%s"}[5m])) by(name)`))
	for _, c := range categories {
		if c == "java" {
			charts = append(charts, model.NewChartInfo("threads", "Threads", "instance", `jvm_threads_current{service="%s"}`))
			charts = append(charts, model.NewChartInfo("gc_duration", "GC Duration", "instance", `rate(jvm_gc_collection_seconds_sum{service="%s"}[1m])`))
		} else if c == "go" {
			charts = append(charts, model.NewChartInfo("threads", "Threads", "instance", `go_threads{service="%s"}`))
			charts = append(charts, model.NewChartInfo("goroutines", "Goroutines", "instance", `go_goroutines{service="%s"}`))
			charts = append(charts, model.NewChartInfo("gc_duration", "GC Duration", "instance", `sum(go_gc_duration_seconds{service="%s"}) by (instance)`))
		}
	}
	for i, c := range charts {
		charts[i].Query = fmt.Sprintf(c.Query, service)
	}
	return
}

func (b *metricBiz) GetMatrix(query, label string, start, end time.Time) (lines []model.ChartLine, err error) {
	api, err := b.getAPI()
	if err != nil {
		return nil, err
	}

	period := end.Sub(start)
	value, err := api.QueryRange(context.Background(), query, papi.Range{
		Start: start,
		End:   end,
		Step:  b.calcStep(period),
	})
	if err != nil {
		return nil, err
	}

	matrix := value.(pmodel.Matrix)
	for _, stream := range matrix {
		line := model.ChartLine{Label: string(stream.Metric[pmodel.LabelName(label)])}
		for _, v := range stream.Values {
			p := model.ChartPoint{
				X: int64(v.Timestamp),
				Y: float64(v.Value),
			}
			line.Data = append(line.Data, p)
		}
		lines = append(lines, line)
	}
	return
}

func (b *metricBiz) calcStep(period time.Duration) (step time.Duration) {
	if period >= times.Day {
		step = 20 * time.Minute
	} else if period >= 12*time.Hour {
		step = 10 * time.Minute
	} else if period >= 6*time.Hour {
		step = 5 * time.Minute
	} else if period >= 3*time.Hour {
		step = 3 * time.Minute
	} else {
		step = time.Minute
	}
	return
}

func (b *metricBiz) getAPI() (api papi.API, err error) {
	v, err := b.api.Get()
	if err != nil {
		return nil, err
	}
	return v.(papi.API), nil
}

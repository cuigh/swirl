package biz

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cuigh/auxo/ext/times"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/auxo/util/lazy"
	"github.com/cuigh/swirl/misc"
	client "github.com/prometheus/client_golang/api"
	papi "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type MetricBiz interface {
	Enabled() bool
	GetMatrix(query, legend string, start, end time.Time) (data *MatrixData, err error)
	GetScalar(query string, t time.Time) (data float64, err error)
	GetVector(query, label string, t time.Time) (data *VectorData, err error)
}

func NewMetric(setting *misc.Setting) MetricBiz {
	b := &metricBiz{prometheus: setting.Metric.Prometheus}
	b.api.New = b.createAPI
	return b
}

type metricBiz struct {
	prometheus string
	api        lazy.Value
}

func (b *metricBiz) createAPI() (api interface{}, err error) {
	if b.prometheus == "" {
		return nil, errors.New("prometheus address is not configured")
	}

	var c client.Client
	if c, err = client.NewClient(client.Config{Address: b.prometheus}); err == nil {
		api = papi.NewAPI(c)
	}
	return
}

func (b *metricBiz) Enabled() bool {
	return b.prometheus != ""
}

func (b *metricBiz) GetMatrix(query, legend string, start, end time.Time) (data *MatrixData, err error) {
	if !b.Enabled() {
		return
	}

	api, err := b.getAPI()
	if err != nil {
		return nil, err
	}

	period := end.Sub(start)
	value, _, err := api.QueryRange(context.Background(), query, papi.Range{
		Start: start,
		End:   end,
		Step:  b.calcStep(period),
	})
	if err != nil {
		return nil, err
	}

	data = &MatrixData{}
	matrix := value.(model.Matrix)
	for _, stream := range matrix {
		data.Legend = append(data.Legend, b.formatLabel(legend, stream.Metric))
		line := MatrixLine{Name: b.formatLabel(legend, stream.Metric)}
		for _, v := range stream.Values {
			p := MatrixPoint{
				X: int64(v.Timestamp),
				Y: float64(v.Value),
			}
			line.Data = append(line.Data, p)
		}
		data.Series = append(data.Series, line)
	}
	return
}

func (b *metricBiz) GetScalar(query string, t time.Time) (v float64, err error) {
	if !b.Enabled() {
		return
	}

	api, err := b.getAPI()
	if err != nil {
		return 0, err
	}

	value, _, err := api.Query(context.Background(), query, t)
	if err != nil {
		return 0, err
	}

	//scalar := value.(*model.Scalar)
	vector := value.(model.Vector)
	if len(vector) > 0 {
		sample := vector[0]
		return float64(sample.Value), nil
	}
	return 0, nil
}

func (b *metricBiz) GetVector(query, label string, t time.Time) (data *VectorData, err error) {
	if !b.Enabled() {
		return
	}

	var api papi.API
	api, err = b.getAPI()
	if err != nil {
		return
	}

	var value model.Value
	value, _, err = api.Query(context.Background(), query, t)
	if err != nil {
		return
	}

	data = &VectorData{}
	vector := value.(model.Vector)
	for _, sample := range vector {
		cv := VectorValue{
			Name:  b.formatLabel(label, sample.Metric),
			Value: float64(sample.Value),
		}
		data.Data = append(data.Data, cv)
		data.Legend = append(data.Legend, cv.Name)
	}
	return
}

func (b *metricBiz) calcStep(period time.Duration) (step time.Duration) {
	if period >= times.Day {
		step = 8 * time.Minute
	} else if period >= 12*time.Hour {
		step = 4 * time.Minute
	} else if period >= 6*time.Hour {
		step = 2 * time.Minute
	} else if period >= 3*time.Hour {
		step = time.Minute
	} else {
		step = 30 * time.Second
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

func (b *metricBiz) formatLabel(label string, metric model.Metric) string {
	return os.Expand(label, func(key string) string {
		if s := string(metric[model.LabelName(key)]); s != "" {
			return s
		}
		return "[" + key + "]"
	})
}

type MatrixData struct {
	Legend []string     `json:"legend"`
	Series []MatrixLine `json:"series"`
}

type MatrixLine struct {
	Name string        `json:"name"`
	Data []MatrixPoint `json:"data"`
}

type MatrixPoint struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
}

func (p *MatrixPoint) MarshalJSON() ([]byte, error) {
	return cast.StringToBytes(fmt.Sprintf("[%v,%v]", p.X, p.Y)), nil
}

type VectorData struct {
	Legend []string      `json:"legend"`
	Data   []VectorValue `json:"data"`
}

type VectorValue struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

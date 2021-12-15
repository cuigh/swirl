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

type ChartBiz interface {
	Search(title, dashboard string, pageIndex, pageSize int) (charts []*Chart, total int, err error)
	Delete(id, title string, user web.User) (err error)
	Find(id string) (chart *Chart, err error)
	Batch(ids ...string) (charts []*model.Chart, err error)
	Create(chart *Chart, user web.User) (err error)
	Update(chart *Chart, user web.User) (err error)
	FetchData(key string, ids []string, period time.Duration) (data.Map, error)
	FindDashboard(name, key string) (dashboard *Dashboard, err error)
	UpdateDashboard(dashboard *model.Dashboard, user web.User) (err error)
}

func NewChart(d dao.Interface, mb MetricBiz, eb EventBiz) ChartBiz {
	return &chartBiz{
		d:  d,
		mb: mb,
		eb: eb,
		builtin: []*model.Chart{
			model.NewChart("service", "$cpu", "CPU", "${name}", `rate(container_cpu_user_seconds_total{container_label_com_docker_swarm_service_name="${service}"}[5m]) * 100`, "percent:100", 60),
			model.NewChart("service", "$memory", "Memory", "${name}", `container_memory_usage_bytes{container_label_com_docker_swarm_service_name="${service}"}`, "size:bytes", 60),
			model.NewChart("service", "$network_in", "Network Receive", "${name}", `sum(irate(container_network_receive_bytes_total{container_label_com_docker_swarm_service_name="${service}"}[5m])) by(name)`, "size:bytes", 60),
			model.NewChart("service", "$network_out", "Network Send", "${name}", `sum(irate(container_network_transmit_bytes_total{container_label_com_docker_swarm_service_name="${service}"}[5m])) by(name)`, "size:bytes", 60),
		},
	}
}

type chartBiz struct {
	builtin []*model.Chart
	d       dao.Interface
	mb      MetricBiz
	eb      EventBiz
}

func (b *chartBiz) Search(title, dashboard string, pageIndex, pageSize int) (charts []*Chart, total int, err error) {
	var list []*model.Chart
	list, total, err = b.d.ChartList(context.TODO(), title, dashboard, pageIndex, pageSize)
	if err == nil {
		charts = make([]*Chart, len(list))
		for i, c := range list {
			charts[i] = newChart(c)
		}
	}
	return
}

func (b *chartBiz) Create(chart *Chart, user web.User) (err error) {
	c := &model.Chart{
		ID:        createId(),
		CreatedAt: time.Now(),
	}
	c.UpdatedAt = c.CreatedAt
	if err = copier.CopyWithOption(c, chart, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	err = b.d.ChartCreate(context.TODO(), c)
	if err == nil {
		b.eb.CreateChart(EventActionCreate, c.ID, c.Title, user)
	}
	return
}

func (b *chartBiz) Delete(id, title string, user web.User) (err error) {
	err = b.d.ChartDelete(context.TODO(), id)
	if err == nil {
		b.eb.CreateChart(EventActionDelete, id, title, user)
	}
	return
}

func (b *chartBiz) Find(id string) (chart *Chart, err error) {
	var c *model.Chart
	c, err = b.d.ChartGet(context.TODO(), id)
	if err == nil {
		chart = newChart(c)
	}
	return
}

func (b *chartBiz) Batch(ids ...string) (charts []*model.Chart, err error) {
	charts, err = b.d.ChartBatch(context.TODO(), ids...)
	return
}

func (b *chartBiz) Update(chart *Chart, user web.User) (err error) {
	c := &model.Chart{
		UpdatedAt: time.Now(),
	}
	if err = copier.CopyWithOption(c, chart, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	err = b.d.ChartUpdate(context.TODO(), c)
	if err == nil {
		b.eb.CreateChart(EventActionUpdate, chart.ID, chart.Title, user)
	}
	return
}

func (b *chartBiz) FindDashboard(name, key string) (dashboard *Dashboard, err error) {
	var d *model.Dashboard
	d, err = b.d.DashboardGet(context.TODO(), name, key)
	if err != nil {
		return
	}

	if d == nil {
		dashboard = defaultDashboard(name, key)
	} else {
		dashboard = newDashboard(d)
	}
	err = b.fillCharts(dashboard)
	return
}

func (b *chartBiz) UpdateDashboard(dashboard *model.Dashboard, user web.User) (err error) {
	return b.d.DashboardUpdate(context.TODO(), dashboard)
}

func (b *chartBiz) FetchData(key string, ids []string, period time.Duration) (data.Map, error) {
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

func (b *chartBiz) fetchMatrixData(chart *model.Chart, key string, start, end time.Time) (md *MatrixData, err error) {
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

func (b *chartBiz) fetchVectorData(chart *model.Chart, key string, end time.Time) (cvd *VectorData, err error) {
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

func (b *chartBiz) fetchScalarData(chart *model.Chart, key string, end time.Time) (*VectorValue, error) {
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

func (b *chartBiz) getCharts(ids []string) (charts map[string]*model.Chart, err error) {
	var (
		customIds    []string
		customCharts []*model.Chart
	)

	charts = make(map[string]*model.Chart)
	for _, id := range ids {
		if id[0] == '$' {
			for _, c := range b.builtin {
				if c.ID == id {
					charts[id] = c
				}
			}
		} else {
			customIds = append(customIds, id)
		}
	}

	if len(customIds) > 0 {
		if customCharts, err = b.Batch(customIds...); err == nil {
			for _, chart := range customCharts {
				charts[chart.ID] = chart
			}
		}
	}
	return
}

func (b *chartBiz) fillCharts(d *Dashboard) (err error) {
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

	for _, info := range d.Charts {
		if c := m[info.ID]; c != nil {
			info.Type = c.Type
			info.Title = c.Title
			info.Unit = c.Unit
			if info.Width == 0 {
				info.Width = c.Width
			}
			if info.Height == 0 {
				info.Height = c.Height
			}
			info.Margin.Left = c.Margin.Left
			info.Margin.Right = c.Margin.Right
			info.Margin.Top = c.Margin.Top
			info.Margin.Bottom = c.Margin.Bottom
		}
	}
	return nil
}

type Dashboard struct {
	Name     string   `json:"name"`
	Key      string   `json:"key,omitempty"`
	Period   int32    `json:"period,omitempty"`   // data range in minutes
	Interval int32    `json:"interval,omitempty"` // refresh interval in seconds, 0 means disabled.
	Charts   []*Chart `json:"charts,omitempty"`
}

type Chart struct {
	ID          string `json:"id"`
	Title       string `json:"title" valid:"required"`
	Description string `json:"desc,omitempty"`
	Dashboard   string `json:"dashboard,omitempty"`
	Type        string `json:"type" valid:"required"`
	Unit        string `json:"unit"`
	Width       int32  `json:"width" valid:"required"`
	Height      int32  `json:"height" valid:"required"`
	Margin      struct {
		Left   int32 `json:"left,omitempty"`
		Right  int32 `json:"right,omitempty"`
		Top    int32 `json:"top,omitempty"`
		Bottom int32 `json:"bottom,omitempty"`
	} `json:"margin"`
	Metrics   []model.ChartMetric `json:"metrics" valid:"required"`
	Options   data.Map            `json:"options,omitempty"`
	CreatedAt string              `json:"createdAt,omitempty" copier:"-"`
	UpdatedAt string              `json:"updatedAt,omitempty" copier:"-"`
}

func newChart(c *model.Chart) *Chart {
	chart := &Chart{
		CreatedAt: formatTime(c.CreatedAt),
		UpdatedAt: formatTime(c.UpdatedAt),
	}
	if err := copier.CopyWithOption(chart, c, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		panic(err)
	}
	return chart
}

func defaultDashboard(name, key string) *Dashboard {
	d := &Dashboard{
		Name:     name,
		Key:      key,
		Period:   30,
		Interval: 15,
	}
	if name == "service" {
		d.Charts = []*Chart{
			{ID: "$cpu"},
			{ID: "$memory"},
			{ID: "$network_in"},
			{ID: "$network_out"},
		}
	}
	return d
}

func newDashboard(d *model.Dashboard) *Dashboard {
	dashboard := &Dashboard{
		Name:     d.Name,
		Key:      d.Key,
		Period:   d.Period,
		Interval: d.Interval,
	}
	if len(d.Charts) > 0 {
		dashboard.Charts = make([]*Chart, len(d.Charts))
		for i, c := range d.Charts {
			dashboard.Charts[i] = &Chart{
				ID:     c.ID,
				Width:  c.Width,
				Height: c.Height,
			}
		}
	}
	return dashboard
}

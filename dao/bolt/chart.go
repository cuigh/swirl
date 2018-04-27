package bolt

import (
	"github.com/cuigh/swirl/model"
)

func (d *Dao) ChartList() (charts []*model.Chart, err error) {
	err = d.each("chart", func(v Value) error {
		chart := &model.Chart{}
		err = v.Unmarshal(chart)
		if err == nil {
			charts = append(charts, chart)
		}
		return err
	})
	return
}

func (d *Dao) ChartCreate(chart *model.Chart) (err error) {
	return d.update("chart", chart.Name, chart)
}

func (d *Dao) ChartGet(name string) (chart *model.Chart, err error) {
	var v Value
	v, err = d.get("chart", name)
	if err == nil {
		if v != nil {
			chart = &model.Chart{}
			err = v.Unmarshal(chart)
		}
	}
	return
}

func (d *Dao) ChartBatch(names ...string) (charts []*model.Chart, err error) {
	err = d.slice("chart", func(v Value) error {
		chart := &model.Chart{}
		err = v.Unmarshal(chart)
		if err == nil {
			charts = append(charts, chart)
		}
		return err
	}, names...)
	return
}

func (d *Dao) ChartUpdate(chart *model.Chart) (err error) {
	return d.update("chart", chart.Name, chart)
}

func (d *Dao) ChartDelete(name string) (err error) {
	return d.delete("chart", name)
}

func (d *Dao) DashboardGet(name, key string) (dashboard *model.ChartDashboard, err error) {
	dashboard = &model.ChartDashboard{
		Name: name,
		Key:  key,
	}

	var v Value
	v, err = d.get("dashboard", dashboard.ID())
	if err == nil {
		if v != nil {
			err = v.Unmarshal(dashboard)
		}
	}
	return
}

func (d *Dao) DashboardUpdate(dashboard *model.ChartDashboard) (err error) {
	return d.update("dashboard", dashboard.ID(), dashboard)
}

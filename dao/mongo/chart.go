package mongo

import (
	"github.com/cuigh/swirl/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (d *Dao) ChartList() (charts []*model.Chart, err error) {
	d.do(func(db *database) {
		charts = []*model.Chart{}
		err = db.C("chart").Find(nil).All(&charts)
	})
	return
}

func (d *Dao) ChartCreate(chart *model.Chart) (err error) {
	d.do(func(db *database) {
		err = db.C("chart").Insert(chart)
	})
	return
}

func (d *Dao) ChartGet(name string) (chart *model.Chart, err error) {
	d.do(func(db *database) {
		chart = &model.Chart{}
		err = db.C("chart").FindId(name).One(chart)
		if err == mgo.ErrNotFound {
			chart, err = nil, nil
		} else if err != nil {
			chart = nil
		}
	})
	return
}

func (d *Dao) ChartBatch(names ...string) (charts []*model.Chart, err error) {
	d.do(func(db *database) {
		q := bson.M{"_id": bson.M{"$in": names}}
		charts = make([]*model.Chart, 0)
		err = db.C("chart").Find(q).All(&charts)
	})
	return
}

func (d *Dao) ChartUpdate(chart *model.Chart) (err error) {
	d.do(func(db *database) {
		err = db.C("chart").UpdateId(chart.Name, chart)
	})
	return
}

func (d *Dao) ChartDelete(name string) (err error) {
	d.do(func(db *database) {
		err = db.C("chart").RemoveId(name)
	})
	return
}

func (d *Dao) DashboardGet(name, key string) (dashboard *model.ChartDashboard, err error) {
	d.do(func(db *database) {
		dashboard = &model.ChartDashboard{
			Name: name,
			Key:  key,
		}
		err = db.C("dashboard").FindId(dashboard.ID()).One(dashboard)
		if err == mgo.ErrNotFound {
			dashboard, err = nil, nil
		} else if err != nil {
			dashboard = nil
		}
	})
	return
}

func (d *Dao) DashboardUpdate(dashboard *model.ChartDashboard) (err error) {
	d.do(func(db *database) {
		update := bson.M{
			"$set": dashboard,
		}
		_, err = db.C("dashboard").UpsertId(dashboard.ID(), update)
	})
	return
}

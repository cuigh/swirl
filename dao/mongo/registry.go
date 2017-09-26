package mongo

import (
	"time"

	"github.com/cuigh/swirl/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (d *Dao) RegistryCreate(registry *model.Registry) (err error) {
	d.do(func(db *database) {
		err = db.C("registry").Insert(registry)
	})
	return
}

func (d *Dao) RegistryUpdate(registry *model.Registry) (err error) {
	d.do(func(db *database) {
		update := bson.M{
			"name":       registry.Name,
			"url":        registry.URL,
			"username":   registry.Username,
			"updated_at": time.Now(),
		}
		if registry.Password != "" {
			update["password"] = registry.Password
		}
		err = db.C("registry").UpdateId(registry.ID, bson.M{"$set": update})
	})
	return
}

func (d *Dao) RegistryList() (registries []*model.Registry, err error) {
	d.do(func(db *database) {
		registries = []*model.Registry{}
		err = db.C("registry").Find(nil).All(&registries)
	})
	return
}

func (d *Dao) RegistryGet(id string) (registry *model.Registry, err error) {
	d.do(func(db *database) {
		registry = &model.Registry{}
		err = db.C("registry").FindId(id).One(registry)
		if err == mgo.ErrNotFound {
			err = nil
		} else if err != nil {
			registry = nil
		}
	})
	return
}

func (d *Dao) RegistryDelete(id string) (err error) {
	d.do(func(db *database) {
		err = db.C("registry").RemoveId(id)
	})
	return
}

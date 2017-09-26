package mongo

import (
	"github.com/cuigh/swirl/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const settingID int32 = 0

func (d *Dao) SettingGet() (setting *model.Setting, err error) {
	d.do(func(db *database) {
		setting = &model.Setting{}
		err = db.C("setting").FindId(settingID).One(setting)
		if err == mgo.ErrNotFound {
			err = nil
		}
	})
	return
}

func (d *Dao) SettingUpdate(setting *model.Setting) (err error) {
	d.do(func(db *database) {
		update := bson.M{
			"$set": setting,
		}
		_, err = db.C("setting").UpsertId(settingID, update)
	})
	return
}

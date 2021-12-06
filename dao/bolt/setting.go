package bolt

import (
	"context"
	"time"

	"github.com/cuigh/swirl/model"
)

func (d *Dao) SettingList(ctx context.Context) (settings []*model.Setting, err error) {
	err = d.each("setting", func(v Value) error {
		s := &model.Setting{}
		err = v.Unmarshal(s)
		if err != nil {
			return err
		}

		settings = append(settings, s)
		return nil
	})
	return
}

func (d *Dao) SettingGet(ctx context.Context, id string) (setting *model.Setting, err error) {
	var v Value
	v, err = d.get("setting", id)
	if err == nil && v != nil {
		setting = &model.Setting{}
		err = v.Unmarshal(setting)
	}
	return
}

func (d *Dao) SettingUpdate(ctx context.Context, id string, options []*model.SettingOption) (err error) {
	setting := &model.Setting{
		ID:        id,
		Options:   options,
		UpdatedAt: time.Now(),
	}
	return d.update("setting", id, setting)
}

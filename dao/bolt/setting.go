package bolt

import (
	"context"

	"github.com/cuigh/swirl/dao"
)

const Setting = "setting"

func (d *Dao) SettingGetAll(ctx context.Context) (settings []*dao.Setting, err error) {
	err = d.each(Setting, func(v []byte) error {
		s := &dao.Setting{}
		err = decode(v, s)
		if err != nil {
			return err
		}

		settings = append(settings, s)
		return nil
	})
	return
}

func (d *Dao) SettingGet(ctx context.Context, id string) (setting *dao.Setting, err error) {
	setting = &dao.Setting{}
	err = d.get(Setting, id, setting)
	if err == ErrNoRecords {
		return nil, nil
	} else if err != nil {
		setting = nil
	}
	return
}

func (d *Dao) SettingUpdate(ctx context.Context, setting *dao.Setting) (err error) {
	return d.replace(Setting, setting.ID, setting)
}

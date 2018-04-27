package bolt

import (
	"github.com/cuigh/swirl/model"
)

const settingID = "0"

func (d *Dao) SettingGet() (setting *model.Setting, err error) {
	var v Value
	v, err = d.get("setting", settingID)
	if err == nil {
		setting = &model.Setting{}
		if v != nil {
			err = v.Unmarshal(setting)
		}
	}
	return
}

func (d *Dao) SettingUpdate(setting *model.Setting) (err error) {
	return d.update("setting", settingID, setting)
}

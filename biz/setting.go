package biz

import (
	"time"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
)

// Setting return a setting biz instance.
var Setting = &settingBiz{}

type settingBiz struct {
}

func (b *settingBiz) Get() (setting *model.Setting, err error) {
	do(func(d dao.Interface) {
		setting, err = d.SettingGet()
	})
	return
}

func (b *settingBiz) Update(setting *model.Setting, user web.User) (err error) {
	do(func(d dao.Interface) {
		setting.UpdatedBy = user.ID()
		setting.UpdatedAt = time.Now()
		err = d.SettingUpdate(setting)
		if err == nil {
			Event.CreateSetting(model.EventActionUpdate, user)
		}
	})
	return
}

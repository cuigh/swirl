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
	loc *time.Location
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

func (b *settingBiz) Time(t time.Time) string {
	if b.loc == nil {
		// todo: auto refresh settings after update
		if s, err := b.Get(); err == nil && s != nil {
			b.loc = time.FixedZone(s.TimeZone.Name, int(s.TimeZone.Offset))
		} else {
			b.loc = time.Local
		}
	}
	return t.In(b.loc).Format("2006-01-02 15:04:05")
}

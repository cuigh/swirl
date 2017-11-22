package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

// SettingController is a controller of system setting
type SettingController struct {
	Index  web.HandlerFunc `path:"/" name:"setting.edit" authorize:"!" desc:"settings edit page"`
	Update web.HandlerFunc `path:"/" name:"setting.update" method:"post" authorize:"!" desc:"update settings"`
}

// Setting creates an instance of SettingController
func Setting() (c *SettingController) {
	return &SettingController{
		Index:  settingIndex,
		Update: settingUpdate,
	}
}

func settingIndex(ctx web.Context) error {
	setting, err := biz.Setting.Get()
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Setting", setting).Set("TimeZones", misc.TimeZones)
	return ctx.Render("system/setting/index", m)
}

func settingUpdate(ctx web.Context) error {
	setting := &model.Setting{}
	err := ctx.Bind(setting)
	if err == nil {
		for _, tz := range misc.TimeZones {
			if tz.Name == setting.TimeZone.Name {
				setting.TimeZone.Offset = tz.Offset
				break
			}
		}
		err = biz.Setting.Update(setting, ctx.User())
	}
	return ajaxResult(ctx, err)
}

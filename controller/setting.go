package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

type SettingController struct {
	Index  web.HandlerFunc `path:"/" name:"setting.edit" authorize:"!" desc:"settings edit page"`
	Update web.HandlerFunc `path:"/" name:"setting.update" method:"post" authorize:"!" desc:"update settings"`
}

func Setting() (c *SettingController) {
	c = &SettingController{}

	c.Index = func(ctx web.Context) error {
		setting, err := biz.Setting.Get()
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Setting", setting)
		return ctx.Render("system/setting/index", m)
	}

	c.Update = func(ctx web.Context) error {
		setting := &model.Setting{}
		err := ctx.Bind(setting)
		if err == nil {
			err = biz.Setting.Update(setting, ctx.User())
		}
		return ajaxResult(ctx, err)
	}

	return
}

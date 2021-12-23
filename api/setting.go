package api

import (
	"encoding/json"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
)

// SettingHandler encapsulates setting related handlers.
type SettingHandler struct {
	Load web.HandlerFunc `path:"/load" auth:"setting.view" desc:"load setting"`
	Save web.HandlerFunc `path:"/save" method:"post" auth:"setting.edit" desc:"save setting"`
}

// NewSetting creates an instance of SettingHandler
func NewSetting(b biz.SettingBiz) *SettingHandler {
	return &SettingHandler{
		Load: settingLoad(b),
		Save: settingSave(b),
	}
}

func settingLoad(b biz.SettingBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		options, err := b.Load()
		if err != nil {
			return err
		}
		return success(ctx, options)
	}
}

func settingSave(b biz.SettingBiz) web.HandlerFunc {
	type Args struct {
		ID      string          `json:"id"`
		Options json.RawMessage `json:"options"`
	}

	return func(ctx web.Context) (err error) {
		args := &Args{}
		err = ctx.Bind(args)
		if err == nil {
			err = b.Save(args.ID, args.Options, ctx.User())
		}
		return ajax(ctx, err)
	}
}

package biz

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
)

type SettingBiz interface {
	Find(ctx context.Context, id string) (options interface{}, err error)
	Load(ctx context.Context) (options data.Map, err error)
	Save(ctx context.Context, id string, options interface{}, user web.User) (err error)
}

func NewSetting(d dao.Interface, eb EventBiz) SettingBiz {
	return &settingBiz{d: d, eb: eb}
}

type settingBiz struct {
	d  dao.Interface
	eb EventBiz
}

func (b *settingBiz) Find(ctx context.Context, id string) (options interface{}, err error) {
	var setting *dao.Setting
	setting, err = b.d.SettingGet(ctx, id)
	if err == nil && setting != nil {
		return b.unmarshal(setting.Options)
	}
	return
}

// Load returns settings of swirl. If not found, default settings will be returned.
func (b *settingBiz) Load(ctx context.Context) (options data.Map, err error) {
	var settings []*dao.Setting
	settings, err = b.d.SettingGetAll(ctx)
	if err != nil {
		return
	}

	options = data.Map{}
	for _, s := range settings {
		var v interface{}
		if v, err = b.unmarshal(s.Options); err != nil {
			return
		}
		options[s.ID] = v
	}
	return
}

func (b *settingBiz) Save(ctx context.Context, id string, options interface{}, user web.User) (err error) {
	setting := &dao.Setting{
		ID:        id,
		UpdatedAt: time.Now(),
	}
	if user != nil {
		setting.UpdatedBy = dao.Operator{ID: user.ID(), Name: user.Name()}
	}

	setting.Options, err = b.marshal(options)
	if err == nil {
		err = b.d.SettingUpdate(ctx, setting)
	}
	if err == nil && user != nil {
		b.eb.CreateSetting(EventActionUpdate, user)
	}
	return
}

func (b *settingBiz) marshal(v interface{}) (s string, err error) {
	var buf []byte
	if buf, err = json.Marshal(v); err == nil {
		s = string(buf)
	}
	return
}

func (b *settingBiz) unmarshal(s string) (v interface{}, err error) {
	d := json.NewDecoder(bytes.NewBuffer([]byte(s)))
	d.UseNumber()
	err = d.Decode(&v)
	return
}

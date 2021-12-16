package biz

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
)

type SettingBiz interface {
	Find(id string) (options data.Map, err error)
	Load() (options data.Map, err error)
	Save(id string, options data.Map, user web.User) (err error)
}

func NewSetting(d dao.Interface, eb EventBiz) SettingBiz {
	return &settingBiz{d: d, eb: eb}
}

type settingBiz struct {
	d  dao.Interface
	eb EventBiz
}

func (b *settingBiz) Find(id string) (options data.Map, err error) {
	var setting *model.Setting
	setting, err = b.d.SettingGet(context.TODO(), id)
	if err != nil {
		return
	}

	if setting != nil {
		options = b.toMap(setting.Options)
	} else {
		options = make(data.Map)
	}
	return
}

// Load returns settings of swirl. If not found, default settings will be returned.
func (b *settingBiz) Load() (options data.Map, err error) {
	var settings []*model.Setting
	settings, err = b.d.SettingGetAll(context.TODO())
	if err != nil {
		return
	}

	options = data.Map{}
	for _, s := range settings {
		options[s.ID] = b.toMap(s.Options)
	}
	return
}

func (b *settingBiz) Save(id string, options data.Map, user web.User) (err error) {
	setting := &model.Setting{
		ID:        id,
		Options:   b.toOptions(options),
		UpdatedAt: time.Now(),
		UpdatedBy: model.Operator{ID: user.ID(), Name: user.Name()},
	}
	err = b.d.SettingUpdate(context.TODO(), setting)
	if err == nil && user != nil {
		b.eb.CreateSetting(EventActionUpdate, user)
	}
	return
}

func (b *settingBiz) toOptions(m data.Map) []*model.SettingOption {
	var opts []*model.SettingOption
	for k, v := range m {
		opt := &model.SettingOption{Name: k}
		switch v.(type) {
		case bool:
			opt.Type = "bool"
			opt.Value = strconv.FormatBool(v.(bool))
		case json.Number:
			opt.Type = "number"
			opt.Value = cast.ToString(v)
		case string:
			opt.Type = "string"
			opt.Value = v.(string)
		default:
			opt.Type = "json"
			opt.Value = b.toJSON(v)
		}
		opts = append(opts, opt)
	}
	return opts
}

func (b *settingBiz) toMap(options []*model.SettingOption) data.Map {
	m := data.Map{}
	for _, opt := range options {
		var v interface{}
		switch opt.Type {
		case "bool":
			v = opt.Value == "true"
		case "number":
			v = cast.ToInt32(opt.Value)
		case "string":
			v = opt.Value
		default:
			v = b.fromJSON(opt.Value)
		}
		m[opt.Name] = v
	}
	return m
}

func (b *settingBiz) toJSON(v interface{}) string {
	d, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(d)
}

func (b *settingBiz) fromJSON(v string) interface{} {
	var i interface{}
	err := json.Unmarshal([]byte(v), &i)
	if err != nil {
		panic(err)
	}
	return i
}

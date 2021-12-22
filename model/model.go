package model

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/ext/times"
	"github.com/docker/docker/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type Time time.Time

func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(time.Time(t))
}

func (t *Time) UnmarshalBSONValue(bt bsontype.Type, data []byte) error {
	if v, _, valid := bsoncore.ReadValue(data, bt); valid {
		*t = Time(v.Time())
		return nil
	}
	return errors.Format("unmarshal failed, type: %s, data:%s", bt, data)
}

func (t Time) MarshalJSON() (b []byte, err error) {
	return strconv.AppendInt(b, times.ToUnixMilli(time.Time(t)), 10), nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		*t = Time(times.FromUnixMilli(i))
	}
	return err
}

func (t Time) String() string {
	return time.Time(t).String()
}

func (t Time) Format(layout string) string {
	return time.Time(t).Format(layout)
}

type Operator struct {
	ID   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

// Setting represents the options of swirl.
type Setting struct {
	ID        string           `json:"id" bson:"_id"`
	Options   []*SettingOption `json:"options" bson:"options"`
	UpdatedAt time.Time        `json:"updatedAt" bson:"updated_at"`
	UpdatedBy Operator         `json:"updatedBy" bson:"updated_by"`
}

type SettingOption struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
	Type  string `json:"type" bson:"type"`
}

type Role struct {
	ID          string   `json:"id,omitempty" bson:"_id"`
	Name        string   `json:"name,omitempty" bson:"name" valid:"required"`
	Description string   `json:"desc,omitempty" bson:"desc,omitempty"`
	Perms       []string `json:"perms,omitempty" bson:"perms,omitempty"`
	UpdatedAt   Time     `json:"updatedAt" bson:"updated_at"`
	CreatedAt   Time     `json:"createdAt" bson:"created_at"`
	CreatedBy   Operator `json:"createdBy" bson:"created_by"`
	UpdatedBy   Operator `json:"updatedBy" bson:"updated_by"`
}

type User struct {
	ID        string   `json:"id,omitempty" bson:"_id"`
	Name      string   `json:"name" bson:"name" valid:"required"`
	LoginName string   `json:"loginName" bson:"login_name" valid:"required"`
	Password  string   `json:"-" bson:"password"`
	Salt      string   `json:"-" bson:"salt"`
	Email     string   `json:"email" bson:"email" valid:"required"`
	Admin     bool     `json:"admin" bson:"admin"`
	Type      string   `json:"type" bson:"type"`
	Status    int32    `json:"status" bson:"status"`
	Roles     []string `json:"roles,omitempty" bson:"roles,omitempty"`
	CreatedAt Time     `json:"createdAt" bson:"created_at"`
	UpdatedAt Time     `json:"updatedAt" bson:"updated_at"`
	CreatedBy Operator `json:"createdBy" bson:"created_by"`
	UpdatedBy Operator `json:"updatedBy" bson:"updated_by"`
}

type UserSearchArgs struct {
	Name      string
	LoginName string
	Admin     bool
	Status    int32
	PageIndex int
	PageSize  int
}

type Registry struct {
	ID        string   `json:"id" bson:"_id"`
	Name      string   `json:"name" bson:"name"`
	URL       string   `json:"url" bson:"url"`
	Username  string   `json:"username" bson:"username"`
	Password  string   `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt Time     `json:"createdAt" bson:"created_at"`
	UpdatedAt Time     `json:"updatedAt" bson:"updated_at"`
	CreatedBy Operator `json:"createdBy" bson:"created_by"`
	UpdatedBy Operator `json:"updatedBy" bson:"updated_by"`
}

func (r *Registry) Match(image string) bool {
	return strings.HasPrefix(image, r.URL)
}

func (r *Registry) GetEncodedAuth() string {
	cfg := &types.AuthConfig{
		ServerAddress: r.URL,
		Username:      r.Username,
		Password:      r.Password,
	}
	if buf, e := json.Marshal(cfg); e == nil {
		return base64.URLEncoding.EncodeToString(buf)
	}
	return ""
}

type Stack struct {
	Name      string   `json:"name" bson:"_id"`
	Content   string   `json:"content" bson:"content"`
	Services  []string `json:"services,omitempty" bson:"-"`
	Internal  bool     `json:"internal" bson:"-"`
	CreatedAt Time     `json:"createdAt" bson:"created_at"`
	UpdatedAt Time     `json:"updatedAt" bson:"updated_at"`
	CreatedBy Operator `json:"createdBy" bson:"created_by"`
	UpdatedBy Operator `json:"updatedBy" bson:"updated_by"`
}

type Event struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Type     string             `json:"type" bson:"type"`
	Action   string             `json:"action" bson:"action"`
	Code     string             `json:"code" bson:"code"`
	Name     string             `json:"name" bson:"name"`
	UserID   string             `json:"userId" bson:"user_id"`
	Username string             `json:"username" bson:"username"`
	Time     Time               `json:"time" bson:"time"`
}

type EventSearchArgs struct {
	Type      string `bind:"type"`
	Name      string `bind:"name"`
	PageIndex int    `bind:"pageIndex"`
	PageSize  int    `bind:"pageSize"`
}

// Chart represents a dashboard chart.
type Chart struct {
	ID          string        `json:"id" bson:"_id"` // the id of built-in charts has '$' prefix.
	Title       string        `json:"title" bson:"title" valid:"required"`
	Description string        `json:"desc" bson:"desc"`
	Metrics     []ChartMetric `json:"metrics" bson:"metrics" valid:"required"`
	Dashboard   string        `json:"dashboard" bson:"dashboard"` // home/service...
	Type        string        `json:"type" bson:"type"`           // pie/line...
	Unit        string        `json:"unit" bson:"unit"`           // bytes/milliseconds/percent:100...
	Width       int32         `json:"width" bson:"width"`         // 1-12(12 columns total)
	Height      int32         `json:"height" bson:"height"`       // default 50
	Options     data.Map      `json:"options,omitempty" bson:"options,omitempty"`
	Margin      struct {
		Left   int32 `json:"left,omitempty" bson:"left,omitempty"`
		Right  int32 `json:"right,omitempty" bson:"right,omitempty"`
		Top    int32 `json:"top,omitempty" bson:"top,omitempty"`
		Bottom int32 `json:"bottom,omitempty" bson:"bottom,omitempty"`
	} `json:"margin" bson:"margin"`
	CreatedAt Time     `json:"createdAt" bson:"created_at"`
	UpdatedAt Time     `json:"updatedAt" bson:"updated_at"`
	CreatedBy Operator `json:"createdBy" bson:"created_by"`
	UpdatedBy Operator `json:"updatedBy" bson:"updated_by"`
}

func NewChart(dashboard, id, title, legend, query, unit string, left int32) *Chart {
	c := &Chart{
		ID:          id,
		Title:       title,
		Description: title,
		Metrics: []ChartMetric{
			{Legend: legend, Query: query},
		},
		Dashboard: dashboard,
		Type:      "line",
		Unit:      unit,
		Width:     12,
		Height:    200,
	}
	c.Margin.Left = left
	return c
}

type ChartMetric struct {
	Legend string `json:"legend"`
	Query  string `json:"query"`
}

type ChartSearchArgs struct {
	Title     string `bind:"title"`
	Dashboard string `bind:"dashboard"`
	PageIndex int    `bind:"pageIndex"`
	PageSize  int    `bind:"pageSize"`
}

type Dashboard struct {
	Name      string      `json:"name" bson:"name"`
	Key       string      `json:"key,omitempty" bson:"key,omitempty"`
	Period    int32       `json:"period,omitempty" bson:"period,omitempty"`     // data range in minutes
	Interval  int32       `json:"interval,omitempty" bson:"interval,omitempty"` // refresh interval in seconds, 0 means disabled.
	Charts    []ChartInfo `json:"charts,omitempty" bson:"charts,omitempty"`
	UpdatedAt Time        `json:"-" bson:"updated_at"`
	UpdatedBy Operator    `json:"-" bson:"updated_by"`
}

type ChartInfo struct {
	ID     string `json:"id" bson:"id"`
	Width  int32  `json:"width,omitempty" bson:"width,omitempty"`
	Height int32  `json:"height,omitempty" bson:"height,omitempty"`
	Title  string `json:"title" bson:"-"`
	Type   string `json:"type" bson:"-"`
	Unit   string `json:"unit" bson:"-"`
	Margin struct {
		Left   int32 `json:"left,omitempty" bson:"-"`
		Right  int32 `json:"right,omitempty" bson:"-"`
		Top    int32 `json:"top,omitempty" bson:"-"`
		Bottom int32 `json:"bottom,omitempty" bson:"-"`
	} `json:"margin" bson:"-"`
}

func (cd *Dashboard) ID() string {
	if cd.Key == "" {
		return cd.Name
	}
	return cd.Name + ":" + cd.Key
}

type Session struct {
	ID        string    `json:"id" bson:"_id"` // token
	UserID    string    `json:"userId" bson:"user_id"`
	Username  string    `json:"username" bson:"username"`
	Admin     bool      `json:"admin" bson:"admin"`
	Roles     []string  `json:"roles" bson:"roles"`
	Perm      uint64    `json:"perm" bson:"perm"`
	Perms     []string  `json:"-" bson:"-"`
	Dirty     bool      `json:"dirty" bson:"dirty"`
	Expiry    time.Time `json:"expiry" bson:"expiry"`
	MaxExpiry time.Time `json:"maxExpiry" bson:"max_expiry"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

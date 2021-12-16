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
	Options   []*SettingOption `json:"options" bson:"options,omitempty"`
	UpdatedAt time.Time        `bson:"updated_at" json:"updatedAt,omitempty"`
	UpdatedBy Operator         `json:"updatedBy" bson:"updated_by"`
}

type SettingOption struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
	Type  string `json:"type" bson:"type"`
}

type Role struct {
	ID          string    `bson:"_id" json:"id,omitempty"`
	Name        string    `bson:"name" json:"name,omitempty" valid:"required"`
	Description string    `bson:"desc" json:"desc,omitempty"`
	Perms       []string  `bson:"perms" json:"perms,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at,omitempty"`
	CreatedBy   Operator  `json:"createdBy" bson:"created_by"`
	UpdatedBy   Operator  `json:"updatedBy" bson:"updated_by"`
}

type User struct {
	ID        string    `bson:"_id" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name,omitempty" valid:"required"`
	LoginName string    `bson:"login_name" json:"loginName,omitempty" valid:"required"`
	Password  string    `bson:"password" json:"-"`
	Salt      string    `bson:"salt" json:"-"`
	Email     string    `bson:"email" json:"email,omitempty" valid:"required"`
	Admin     bool      `bson:"admin" json:"admin,omitempty"`
	Type      string    `bson:"type" json:"type,omitempty"`
	Status    int32     `bson:"status" json:"status,omitempty"`
	Roles     []string  `bson:"roles" json:"roles,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt,omitempty"`
	CreatedBy Operator  `json:"createdBy" bson:"created_by"`
	UpdatedBy Operator  `json:"updatedBy" bson:"updated_by"`
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
	ID        string    `bson:"_id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	URL       string    `bson:"url" json:"url"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"-"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt,omitempty"`
	CreatedBy Operator  `json:"createdBy" bson:"created_by"`
	UpdatedBy Operator  `json:"updatedBy" bson:"updated_by"`
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
	Name      string    `bson:"_id" json:"name,omitempty"`
	Content   string    `bson:"content" json:"content,omitempty" bind:"content=form,content=file"`
	Services  []string  `bson:"-" json:"services,omitempty"`
	Internal  bool      `bson:"-" json:"internal"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt,omitempty"`
	CreatedBy Operator  `json:"createdBy" bson:"created_by"`
	UpdatedBy Operator  `json:"updatedBy" bson:"updated_by"`
}

type Event struct {
	ID       primitive.ObjectID `bson:"_id"`
	Type     string             `bson:"type"`
	Action   string             `bson:"action"`
	Code     string             `bson:"code"`
	Name     string             `bson:"name"`
	UserID   string             `bson:"user_id"`
	Username string             `bson:"username"`
	Time     time.Time          `bson:"time"`
}

type EventSearchArgs struct {
	Type      string `bind:"type"`
	Name      string `bind:"name"`
	PageIndex int    `bind:"pageIndex"`
	PageSize  int    `bind:"pageSize"`
}

// Chart represents a dashboard chart.
type Chart struct {
	ID          string        `json:"id" bson:"_id"` // unique, the name of build-in charts has '$' prefix.
	Title       string        `json:"title" valid:"required"`
	Description string        `json:"desc" bson:"desc"`
	Metrics     []ChartMetric `json:"metrics" valid:"required"`
	Dashboard   string        `json:"dashboard"` // home/service/task...
	Type        string        `json:"type"`      // pie/line...
	Unit        string        `json:"unit"`      // bytes/milliseconds/percent:100...
	Width       int32         `json:"width"`     // 1-12(12 columns total)
	Height      int32         `json:"height"`    // default 50
	Options     data.Map      `json:"options,omitempty"`
	Margin      struct {
		Left   int32 `json:"left,omitempty"`
		Right  int32 `json:"right,omitempty"`
		Top    int32 `json:"top,omitempty"`
		Bottom int32 `json:"bottom,omitempty"`
	} `json:"margin"`
	//Colors      []string `json:"colors"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt,omitempty"`
	CreatedBy Operator  `json:"createdBy" bson:"created_by"`
	UpdatedBy Operator  `json:"updatedBy" bson:"updated_by"`
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
	Name      string      `json:"name"`
	Key       string      `json:"key,omitempty"`
	Period    int32       `json:"period,omitempty"`   // data range in minutes
	Interval  int32       `json:"interval,omitempty"` // refresh interval in seconds, 0 means disabled.
	Charts    []ChartInfo `json:"charts,omitempty"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updatedAt,omitempty"`
	UpdatedBy Operator    `json:"updatedBy" bson:"updated_by"`
}

type ChartInfo struct {
	ID     string `json:"id"`
	Width  int32  `json:"width,omitempty"`
	Height int32  `json:"height,omitempty"`
	//Colors []string `json:"colors,omitempty"`
}

func (cd *Dashboard) ID() string {
	if cd.Key == "" {
		return cd.Name
	}
	return cd.Name + ":" + cd.Key
}

type Session struct {
	UserID    string    `bson:"_id" json:"id,omitempty"`
	Token     string    `bson:"token" json:"token,omitempty"`
	Expires   time.Time `bson:"expires" json:"expires,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
}

//type AuthUser struct {
//	user  *User
//	roles []*Role
//	perms map[string]struct{}
//}
//
//func NewAuthUser(user *User, roles []*Role) *AuthUser {
//	if user == nil {
//		panic(111)
//	}
//	u := &AuthUser{
//		user:  user,
//		roles: roles,
//		perms: make(map[string]struct{}),
//	}
//	for _, role := range roles {
//		for _, perm := range role.Perms {
//			u.perms[perm] = data.Empty
//		}
//	}
//	return u
//}
//
//func (u *AuthUser) ID() string {
//	return u.user.ID
//}
//
//func (u *AuthUser) Name() string {
//	return u.user.Name
//}
//
//func (u *AuthUser) Anonymous() bool {
//	return u.user.ID == ""
//}
//
//func (u *AuthUser) Admin() bool {
//	return u.user.Admin
//}
//
//func (u *AuthUser) IsInRole(roleID string) bool {
//	for _, role := range u.roles {
//		if role.ID == roleID {
//			return true
//		}
//	}
//	return false
//}
//
//func (u *AuthUser) IsAllowed(perm string) bool {
//	if u.user.Admin {
//		return true
//	}
//
//	_, ok := u.perms[perm]
//	return ok
//}

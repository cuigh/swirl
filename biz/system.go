package biz

import (
	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
	"github.com/docker/docker/api/types/versions"
)

type SystemBiz interface {
	Init() (err error)
	CheckState() (state *SystemState, err error)
	CreateAdmin(user *User) (err error)
}

func NewSystem(d dao.Interface, ub UserBiz, sb SettingBiz, s *misc.Setting) SystemBiz {
	return &systemBiz{
		d:  d,
		ub: ub,
		sb: sb,
		s:  s,
	}
}

type systemBiz struct {
	d  dao.Interface
	ub UserBiz
	sb SettingBiz
	s  *misc.Setting
}

func (b *systemBiz) Init() (err error) {
	if versions.LessThan(b.s.System.Version, app.Version) {
		// initialize database
		err = b.d.Init()
		if err == nil {
			err = b.sb.Save("system", data.Map{"version": app.Version}, nil)
		}
	}
	return
}

func (b *systemBiz) CheckState() (state *SystemState, err error) {
	var count int
	count, err = b.ub.Count()
	if err == nil {
		state = &SystemState{Fresh: count == 0}
	}
	return
}

func (b *systemBiz) CreateAdmin(user *User) (err error) {
	user.Admin = true
	user.Type = UserTypeInternal

	var count int
	if count, err = b.ub.Count(); err == nil {
		if count > 0 {
			err = errors.Coded(1, "system was already initialized")
		} else {
			_, err = b.ub.Create(user, nil)
		}
	}
	return
}

type SystemState struct {
	Fresh bool `json:"fresh"`
}

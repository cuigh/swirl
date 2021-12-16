package biz

import (
	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
	"github.com/docker/docker/api/types/versions"
)

type SystemBiz interface {
	Init() (err error)
	CheckState() (state *SystemState, err error)
}

func NewSystem(d dao.Interface, ub UserBiz, sb SettingBiz, s *misc.Setting) SystemBiz {
	return &systemBiz{
		s:  s,
		d:  d,
		ub: ub,
		sb: sb,
	}
}

type systemBiz struct {
	s  *misc.Setting
	d  dao.Interface
	ub UserBiz
	sb SettingBiz
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

type SystemState struct {
	Fresh bool `json:"fresh"`
}

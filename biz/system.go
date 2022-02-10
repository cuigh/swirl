package biz

import (
	"context"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
	"github.com/docker/docker/api/types/versions"
)

type SystemBiz interface {
	Init(ctx context.Context) (err error)
	CheckState(ctx context.Context) (state *SystemState, err error)
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

func (b *systemBiz) Init(ctx context.Context) (err error) {
	if versions.LessThan(b.s.System.Version, app.Version) {
		// upgrade database
		err = b.d.Upgrade(ctx)
		if err == nil {
			err = b.sb.Save(ctx, "system", data.Map{"version": app.Version}, nil)
		}
	}
	return
}

func (b *systemBiz) CheckState(ctx context.Context) (state *SystemState, err error) {
	var count int
	count, err = b.ub.Count(ctx)
	if err == nil {
		state = &SystemState{Fresh: count == 0}
	}
	return
}

type SystemState struct {
	Fresh bool `json:"fresh"`
}

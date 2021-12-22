package biz

import (
	"context"
	"time"

	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
)

type SessionBiz interface {
	Find(token string) (session *model.Session, err error)
	Create(session *model.Session) (err error)
	Update(session *model.Session) (err error)
	UpdateExpiry(id string, expiry time.Time) (err error)
}

func NewSession(d dao.Interface, rb RoleBiz) SessionBiz {
	return &sessionBiz{d: d, rb: rb}
}

type sessionBiz struct {
	d  dao.Interface
	rb RoleBiz
}

func (b *sessionBiz) Find(token string) (session *model.Session, err error) {
	return b.d.SessionGet(context.TODO(), token)
}

func (b *sessionBiz) Create(session *model.Session) (err error) {
	session.CreatedAt = time.Now()
	session.UpdatedAt = session.CreatedAt
	return b.d.SessionCreate(context.TODO(), session)
}

func (b *sessionBiz) Update(session *model.Session) (err error) {
	session.Dirty = false
	session.UpdatedAt = time.Now()
	return b.d.SessionUpdate(context.TODO(), session)
}

func (b *sessionBiz) UpdateExpiry(id string, expiry time.Time) (err error) {
	return b.d.SessionUpdateExpiry(context.TODO(), id, expiry)
}

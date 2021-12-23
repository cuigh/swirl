package biz

import (
	"context"
	"time"

	"github.com/cuigh/swirl/dao"
)

type SessionBiz interface {
	Find(token string) (session *dao.Session, err error)
	Create(session *dao.Session) (err error)
	Update(session *dao.Session) (err error)
	UpdateExpiry(id string, expiry time.Time) (err error)
}

func NewSession(d dao.Interface, rb RoleBiz) SessionBiz {
	return &sessionBiz{d: d, rb: rb}
}

type sessionBiz struct {
	d  dao.Interface
	rb RoleBiz
}

func (b *sessionBiz) Find(token string) (session *dao.Session, err error) {
	return b.d.SessionGet(context.TODO(), token)
}

func (b *sessionBiz) Create(session *dao.Session) (err error) {
	session.CreatedAt = time.Now()
	session.UpdatedAt = session.CreatedAt
	return b.d.SessionCreate(context.TODO(), session)
}

func (b *sessionBiz) Update(session *dao.Session) (err error) {
	session.Dirty = false
	session.UpdatedAt = time.Now()
	return b.d.SessionUpdate(context.TODO(), session)
}

func (b *sessionBiz) UpdateExpiry(id string, expiry time.Time) (err error) {
	return b.d.SessionUpdateExpiry(context.TODO(), id, expiry)
}

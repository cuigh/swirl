package biz

import (
	"context"
	"time"

	"github.com/cuigh/swirl/dao"
)

type SessionBiz interface {
	Find(ctx context.Context, token string) (session *dao.Session, err error)
	Create(ctx context.Context, session *dao.Session) (err error)
	Update(ctx context.Context, session *dao.Session) (err error)
	UpdateExpiry(ctx context.Context, id string, expiry time.Time) (err error)
}

func NewSession(d dao.Interface, rb RoleBiz) SessionBiz {
	return &sessionBiz{d: d, rb: rb}
}

type sessionBiz struct {
	d  dao.Interface
	rb RoleBiz
}

func (b *sessionBiz) Find(ctx context.Context, token string) (session *dao.Session, err error) {
	return b.d.SessionGet(ctx, token)
}

func (b *sessionBiz) Create(ctx context.Context, session *dao.Session) (err error) {
	session.CreatedAt = time.Now()
	session.UpdatedAt = session.CreatedAt
	return b.d.SessionCreate(ctx, session)
}

func (b *sessionBiz) Update(ctx context.Context, session *dao.Session) (err error) {
	session.Dirty = false
	session.UpdatedAt = time.Now()
	return b.d.SessionUpdate(ctx, session)
}

func (b *sessionBiz) UpdateExpiry(ctx context.Context, id string, expiry time.Time) (err error) {
	return b.d.SessionUpdateExpiry(ctx, id, expiry)
}

package biz

import (
	"net/http"
	"strings"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
)

// Perm return a perm biz instance.
var (
	Perm         = &permBiz{}
	ErrForbidden = web.NewError(http.StatusForbidden)
)

type permBiz struct {
}

func (b *permBiz) Delete(resType, resID string, user web.User) (err error) {
	do(func(d dao.Interface) {
		err = d.PermDelete(resType, resID)
	})
	return
}

func (b *permBiz) Get(resType, resID string) (perm *model.Perm, err error) {
	do(func(d dao.Interface) {
		perm, err = d.PermGet(resType, resID)
	})
	return
}

func (b *permBiz) Update(perm *model.Perm, user web.User) (err error) {
	do(func(d dao.Interface) {
		err = d.PermUpdate(perm)
	})
	return
}

func (b *permBiz) Check(user web.User, scope string, resType, resID string) (err error) {
	au := user.(*model.AuthUser)
	if au.Admin() {
		return
	}

	do(func(d dao.Interface) {
		var perm *model.Perm
		perm, err = d.PermGet(resType, resID)
		if err != nil {
			return
		}

		if perm == nil || perm.Scope == model.PermNone || (scope == "read" && perm.Scope == model.PermWrite) {
			return
		}

		for _, u := range perm.Users {
			if user.ID() == u {
				return
			}
		}

		for _, r := range perm.Roles {
			if au.IsInRole(r) {
				return
			}
		}
		err = ErrForbidden
	})
	return
}

func (b *permBiz) Apply(next web.HandlerFunc) web.HandlerFunc {
	return func(ctx web.Context) error {
		opt := ctx.Handler().Option("perm")
		if opt != "" {
			array := strings.Split(opt, ",")
			err := b.Check(ctx.User(), array[0], array[1], ctx.P(array[2]))
			if err != nil {
				return err
			}
		}
		return next(ctx)
	}
}

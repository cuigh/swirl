package controller

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

func permEdit(ctx web.Context, resType, resID, tpl string, m data.Map) error {
	perm, err := biz.Perm.Get(resType, resID)
	if err != nil {
		return err
	}
	if perm == nil {
		perm = &model.Perm{}
	}

	roles, err := biz.Role.List()
	if err != nil {
		return err
	}

	checkedRoles := data.Set{}
	checkedRoles.AddSlice(perm.Roles, func(i int) interface{} {
		return perm.Roles[i]
	})

	var users []*model.User
	for _, id := range perm.Users {
		var user *model.User
		if user, err = biz.User.GetByID(id); err != nil {
			return err
		} else if user != nil {
			users = append(users, user)
		}
	}

	m.Set("Perm", perm).Set("Roles", roles).Set("CheckedRoles", checkedRoles).Set("Users", users)
	return ctx.Render(tpl, m)
}

func permUpdate(resType, argName string) web.HandlerFunc {
	return func(ctx web.Context) error {
		perm := &model.Perm{
			ResType: resType,
			ResID:   ctx.P(argName),
		}
		err := ctx.Bind(perm)
		if err != nil {
			return err
		}

		err = biz.Perm.Update(perm, ctx.User())
		return ajaxResult(ctx, err)
	}
}

package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/security"
)

// UserHandler encapsulates user related handlers.
type UserHandler struct {
	SignIn         web.HandlerFunc `path:"/sign-in" method:"post" auth:"*" desc:"user sign in"`
	Search         web.HandlerFunc `path:"/search" auth:"user.view" desc:"search users"`
	Save           web.HandlerFunc `path:"/save" method:"post" auth:"user.edit" desc:"create or update user"`
	Find           web.HandlerFunc `path:"/find" auth:"user.view" desc:"find user by id"`
	Delete         web.HandlerFunc `path:"/delete" method:"post" auth:"user.delete" desc:"delete user"`
	SetStatus      web.HandlerFunc `path:"/set-status" method:"post" auth:"user.edit" desc:"set user status"`
	ModifyPassword web.HandlerFunc `path:"/modify-password" method:"post" auth:"?" desc:"modify password"`
	ModifyProfile  web.HandlerFunc `path:"/modify-profile" method:"post" auth:"?" desc:"modify profile"`
}

// NewUser creates an instance of UserHandler
func NewUser(b biz.UserBiz, eb biz.EventBiz, auth *security.Identifier) *UserHandler {
	return &UserHandler{
		SignIn:         userSignIn(auth, eb),
		Search:         userSearch(b),
		Save:           userSave(b),
		Find:           userFind(b),
		Delete:         userDelete(b),
		SetStatus:      userSetStatus(b),
		ModifyPassword: userModifyPassword(b),
		ModifyProfile:  userModifyProfile(b),
	}
}

func userSignIn(auth *security.Identifier, eb biz.EventBiz) web.HandlerFunc {
	type SignInArgs struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	return func(c web.Context) (err error) {
		var (
			args = &SignInArgs{}
			user security.Identity
		)

		if err = c.Bind(args); err != nil {
			return err
		}

		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		if user, err = auth.Identify(ctx, args.Name, args.Password); err != nil {
			return err
		}

		eb.CreateUser(biz.EventActionLogin, user.ID(), user.Name(), user)

		return success(c, data.Map{
			"name":  user.Name(),
			"token": user.Token(),
			"perms": user.Perms(),
		})
	}
}

func userSave(b biz.UserBiz) web.HandlerFunc {
	return func(c web.Context) error {
		user := &dao.User{}
		err := c.Bind(user, true)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			if user.ID == "" {
				_, err = b.Create(ctx, user, c.User())
			} else {
				err = b.Update(ctx, user, c.User())
			}
		}
		return ajax(c, err)
	}
}

func userSearch(b biz.UserBiz) web.HandlerFunc {
	type Args struct {
		Filter    string `bind:"filter"` // admins, active, blocked
		Name      string `bind:"name"`
		LoginName string `bind:"loginName"`
		PageIndex int    `bind:"pageIndex"`
		PageSize  int    `bind:"pageSize"`
	}

	return func(c web.Context) error {
		args := &Args{}
		err := c.Bind(args)
		if err != nil {
			return err
		}

		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		users, total, err := b.Search(ctx, args.Name, args.LoginName, args.Filter, args.PageIndex, args.PageSize)
		if err != nil {
			return err
		}
		return success(c, data.Map{"items": users, "total": total})
	}
}

func userFind(b biz.UserBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		id := c.Query("id")
		if id == "" {
			id = c.User().ID()
		}
		user, err := b.FindByID(ctx, id)
		if err != nil {
			return err
		}
		return success(c, user)
	}
}

func userDelete(b biz.UserBiz) web.HandlerFunc {
	type Args struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	return func(c web.Context) error {
		args := &Args{}
		err := c.Bind(args)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Delete(ctx, args.ID, args.Name, c.User())
		}
		return ajax(c, err)
	}
}

func userSetStatus(b biz.UserBiz) web.HandlerFunc {
	type Args struct {
		ID     string `json:"id"`
		Status int32  `json:"status"`
	}

	return func(c web.Context) error {
		args := &Args{}
		err := c.Bind(args)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.SetStatus(ctx, args.ID, args.Status, c.User())
		}
		return ajax(c, err)
	}
}

func userModifyPassword(b biz.UserBiz) web.HandlerFunc {
	type Args struct {
		OldPassword string `json:"oldPwd"`
		NewPassword string `json:"newPwd"`
	}

	return func(c web.Context) error {
		args := &Args{}
		err := c.Bind(args)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.ModifyPassword(ctx, args.OldPassword, args.NewPassword, c.User())
		}
		return ajax(c, err)
	}
}

func userModifyProfile(b biz.UserBiz) web.HandlerFunc {
	return func(c web.Context) error {
		u := &dao.User{}
		err := c.Bind(u, true)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.ModifyProfile(ctx, u, c.User())
		}
		return ajax(c, err)
	}
}

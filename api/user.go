package api

import (
	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
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
func NewUser(b biz.UserBiz, eb biz.EventBiz, auth *security.Authenticator) *UserHandler {
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

func userSignIn(auth *security.Authenticator, eb biz.EventBiz) web.HandlerFunc {
	type SignInArgs struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	return func(ctx web.Context) (err error) {
		var (
			args  = &SignInArgs{}
			user  web.User
			token string
		)

		if err = ctx.Bind(args); err == nil {
			if user, err = auth.Login(args.Name, args.Password); err == nil {
				jwt := container.Find("identifier").(*security.JWT)
				token, err = jwt.CreateToken(user.ID(), user.Name())
			}
		}

		if err != nil {
			return err
		}

		eb.CreateUser(biz.EventActionLogin, user.ID(), user.Name(), user)
		return success(ctx, data.Map{
			"token": token,
			"id":    user.ID(),
			"name":  user.Name(),
		})
	}
}

func userSave(b biz.UserBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		user := &model.User{}
		err := ctx.Bind(user, true)
		if err == nil {
			if user.ID == "" {
				_, err = b.Create(user, ctx.User())
			} else {
				err = b.Update(user, ctx.User())
			}
		}
		return ajax(ctx, err)
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

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err != nil {
			return err
		}

		users, total, err := b.Search(args.Name, args.LoginName, args.Filter, args.PageIndex, args.PageSize)
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"items": users, "total": total})
	}
}

func userFind(b biz.UserBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		id := ctx.Query("id")
		if id == "" {
			id = ctx.User().ID()
		}
		user, err := b.FindByID(id)
		if err != nil {
			return err
		}
		return success(ctx, user)
	}
}

func userDelete(b biz.UserBiz) web.HandlerFunc {
	type Args struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err == nil {
			err = b.Delete(args.ID, args.Name, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func userSetStatus(b biz.UserBiz) web.HandlerFunc {
	type Args struct {
		ID     string `json:"id"`
		Status int32  `json:"status"`
	}

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err == nil {
			err = b.SetStatus(args.ID, args.Status, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func userModifyPassword(b biz.UserBiz) web.HandlerFunc {
	type Args struct {
		OldPassword string `json:"oldPwd"`
		NewPassword string `json:"newPwd"`
	}

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err == nil {
			err = b.ModifyPassword(args.OldPassword, args.NewPassword, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func userModifyProfile(b biz.UserBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		u := &model.User{}
		err := ctx.Bind(u, true)
		if err == nil {
			err = b.ModifyProfile(u, ctx.User())
		}
		return ajax(ctx, err)
	}
}

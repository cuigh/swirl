package biz

import (
	"time"

	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/security/passwd"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

var User = &userBiz{} // nolint: golint

type userBiz struct {
}

func (b *userBiz) GetByID(id string) (user *model.User, err error) {
	do(func(d dao.Interface) {
		user, err = d.UserGetByID(id)
	})
	return
}

func (b *userBiz) GetByName(loginName string) (user *model.User, err error) {
	do(func(d dao.Interface) {
		user, err = d.UserGetByName(loginName)
	})
	return
}

func (b *userBiz) Create(user *model.User, ctxUser web.User) (err error) {
	user.ID = misc.NewID()
	user.Status = model.UserStatusActive
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	if user.Type == model.UserTypeInternal {
		user.Password, user.Salt, err = passwd.Generate(user.Password)
		if err != nil {
			return
		}
	}

	do(func(d dao.Interface) {
		if err = d.UserCreate(user); err == nil && ctxUser != nil {
			Event.CreateUser(model.EventActionCreate, user.LoginName, user.Name, ctxUser)
		}
	})
	return
}

func (b *userBiz) Update(user *model.User, ctxUser web.User) (err error) {
	do(func(d dao.Interface) {
		user.UpdatedAt = time.Now()
		if err = d.UserUpdate(user); err == nil {
			Event.CreateUser(model.EventActionUpdate, user.LoginName, user.Name, ctxUser)
		}
	})
	return
}

func (b *userBiz) Block(id string) (err error) {
	do(func(d dao.Interface) {
		err = d.UserBlock(id, true)
	})
	return
}

func (b *userBiz) Unblock(id string) (err error) {
	do(func(d dao.Interface) {
		err = d.UserBlock(id, false)
	})
	return
}

func (b *userBiz) Delete(id string) (err error) {
	do(func(d dao.Interface) {
		err = d.UserDelete(id)
	})
	return
}

func (b *userBiz) UpdateInfo(user *model.User) (err error) {
	do(func(d dao.Interface) {
		err = d.ProfileUpdateInfo(user)
	})
	return
}

func (b *userBiz) UpdatePassword(id, oldPwd, newPwd string) (err error) {
	do(func(d dao.Interface) {
		var (
			user      *model.User
			pwd, salt string
		)

		user, err = d.UserGetByID(id)
		if err != nil {
			return
		}

		if !passwd.Validate(oldPwd, user.Password, user.Salt) {
			err = errors.New("Current password is incorrect")
			return
		}

		pwd, salt, err = passwd.Generate(newPwd)
		if err != nil {
			return
		}

		err = d.ProfileUpdatePassword(id, pwd, salt)
	})
	return
}

func (b *userBiz) List(args *model.UserListArgs) (users []*model.User, count int, err error) {
	do(func(d dao.Interface) {
		users, count, err = d.UserList(args)
	})
	return
}

func (b *userBiz) Count() (count int, err error) {
	do(func(d dao.Interface) {
		count, err = d.UserCount()
	})
	return
}

func (b *userBiz) UpdateSession(id string) (token string, err error) {
	session := &model.Session{
		UserID:    id,
		Token:     misc.NewID(),
		UpdatedAt: time.Now(),
	}
	session.Expires = session.UpdatedAt.Add(time.Hour * 24)
	do(func(d dao.Interface) {
		err = d.SessionUpdate(session)
		if err == nil {
			token = session.Token
		}
	})
	return
}

func (b *userBiz) GetSession(token string) (session *model.Session, err error) {
	do(func(d dao.Interface) {
		session, err = d.SessionGet(token)
	})
	return
}

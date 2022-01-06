package biz

import (
	"context"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/security/passwd"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
)

const (
	// UserTypeInternal is internal user of swirl
	UserTypeInternal = "internal"
	// UserTypeLDAP is external user of LDAP
	UserTypeLDAP = "ldap"
)

const (
	// UserStatusBlocked is the status which user is blocked
	UserStatusBlocked = 0
	// UserStatusActive is the normal status
	UserStatusActive = 1
)

type UserBiz interface {
	Search(ctx context.Context, name, loginName, filter string, pageIndex, pageSize int) (users []*dao.User, total int, err error)
	Create(ctx context.Context, user *dao.User, ctxUser web.User) (id string, err error)
	Update(ctx context.Context, user *dao.User, ctxUser web.User) (err error)
	FindByID(ctx context.Context, id string) (user *dao.User, err error)
	FindByName(ctx context.Context, loginName string) (user *dao.User, err error)
	FindByToken(ctx context.Context, token string) (user *dao.User, err error)
	FindPrivacy(ctx context.Context, loginName string) (privacy *UserPrivacy, err error)
	Count(ctx context.Context) (count int, err error)
	Delete(ctx context.Context, id, name string, user web.User) (err error)
	SetStatus(ctx context.Context, id string, status int32, user web.User) (err error)
	ModifyPassword(ctx context.Context, oldPwd, newPwd string, user web.User) (err error)
	ModifyProfile(ctx context.Context, user *dao.User, ctxUser web.User) (err error)
}

func NewUser(d dao.Interface, eb EventBiz) UserBiz {
	return &userBiz{d: d, eb: eb}
}

type userBiz struct {
	d  dao.Interface
	eb EventBiz
}

func (b *userBiz) Search(ctx context.Context, name, loginName, filter string, pageIndex, pageSize int) (users []*dao.User, total int, err error) {
	var args = &dao.UserSearchArgs{
		Name:      name,
		LoginName: loginName,
		Status:    -1,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}

	switch filter {
	case "admins":
		args.Admin = true
	case "active":
		args.Status = UserStatusActive
	case "blocked":
		args.Status = UserStatusBlocked
	}

	return b.d.UserSearch(ctx, args)
}

func (b *userBiz) FindByID(ctx context.Context, id string) (user *dao.User, err error) {
	return b.d.UserGet(ctx, id)
}

func (b *userBiz) FindByName(ctx context.Context, loginName string) (user *dao.User, err error) {
	return b.d.UserGetByName(ctx, loginName)
}

func (b *userBiz) FindByToken(ctx context.Context, token string) (user *dao.User, err error) {
	return b.d.UserGetByToken(ctx, token)
}

func (b *userBiz) FindPrivacy(ctx context.Context, loginName string) (privacy *UserPrivacy, err error) {
	var u *dao.User
	u, err = b.d.UserGetByName(ctx, loginName)
	if u != nil {
		privacy = &UserPrivacy{
			ID:       u.ID,
			Name:     u.Name,
			Password: u.Password,
			Salt:     u.Salt,
			Type:     u.Type,
			Status:   u.Status,
		}
	}
	return
}

func (b *userBiz) Create(ctx context.Context, user *dao.User, ctxUser web.User) (id string, err error) {
	user.Tokens = data.Options{data.Option{Name: "test", Value: "abc123"}}
	user.ID = createId()
	user.Status = UserStatusActive
	user.CreatedAt = now()
	if ctxUser != nil {
		user.CreatedBy = newOperator(ctxUser)
	}
	user.UpdatedAt = user.CreatedAt
	user.UpdatedBy = user.CreatedBy
	if user.Type == UserTypeInternal {
		user.Password, user.Salt, err = passwd.Generate(user.Password)
		if err != nil {
			return
		}
	}

	if err = b.d.UserCreate(ctx, user); err == nil && ctxUser != nil {
		b.eb.CreateUser(EventActionCreate, user.LoginName, user.Name, ctxUser)
	}
	id = user.ID
	return
}

func (b *userBiz) Update(ctx context.Context, user *dao.User, ctxUser web.User) (err error) {
	user.UpdatedAt = now()
	user.UpdatedBy = newOperator(ctxUser)
	if err = b.d.UserUpdate(ctx, user); err == nil {
		go func() {
			_ = b.d.SessionUpdateDirty(ctx, user.ID, "")
			b.eb.CreateUser(EventActionUpdate, user.LoginName, user.Name, ctxUser)
		}()
	}
	return
}

func (b *userBiz) SetStatus(ctx context.Context, id string, status int32, user web.User) (err error) {
	u := &dao.User{
		ID:        id,
		Status:    status,
		UpdatedAt: now(),
		UpdatedBy: newOperator(user),
	}
	return b.d.UserUpdateStatus(ctx, u)
}

func (b *userBiz) Delete(ctx context.Context, id, name string, user web.User) (err error) {
	err = b.d.UserDelete(ctx, id)
	if err == nil {
		b.eb.CreateUser(EventActionDelete, id, name, user)
	}
	return
}

func (b *userBiz) ModifyPassword(ctx context.Context, oldPwd, newPwd string, user web.User) (err error) {
	var u *dao.User
	u, err = b.d.UserGet(ctx, user.ID())
	if err != nil {
		return err
	} else if u == nil {
		return errors.Format("user not found: %s", user.ID())
	}

	if !passwd.Validate(oldPwd, u.Password, u.Salt) {
		return errors.Coded(misc.ErrOldPasswordIncorrect, "current password is incorrect")
	}

	if u.Password, u.Salt, err = passwd.Generate(newPwd); err != nil {
		return
	}

	u.UpdatedAt = now()
	u.UpdatedBy = newOperator(user)
	return b.d.UserUpdatePassword(ctx, u)
}

func (b *userBiz) ModifyProfile(ctx context.Context, u *dao.User, user web.User) (err error) {
	u.ID = user.ID()
	u.UpdatedAt = now()
	u.UpdatedBy = newOperator(user)
	return b.d.UserUpdateProfile(ctx, u)
}

func (b *userBiz) Count(ctx context.Context) (count int, err error) {
	return b.d.UserCount(ctx)
}

type UserPrivacy struct {
	ID       string
	Name     string
	Password string `json:"-"`
	Salt     string `json:"-"`
	Type     string
	Status   int32
}

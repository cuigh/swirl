package biz

import (
	"context"
	"time"

	"github.com/cuigh/auxo/data/guid"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/security/passwd"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
	"github.com/jinzhu/copier"
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
	Search(name, loginName, filter string, pageIndex, pageSize int) (users []*User, total int, err error)
	Create(user *User, ctxUser web.User) (id string, err error)
	Update(user *User, ctxUser web.User) (err error)
	FindByID(id string) (user *User, err error)
	FindByName(loginName string) (user *User, err error)
	FindPrivacy(loginName string) (privacy *UserPrivacy, err error)
	Count() (count int, err error)
	Delete(id, name string, user web.User) (err error)
	SetStatus(id string, status int32, user web.User) (err error)
	ModifyPassword(oldPwd, newPwd string, user web.User) (err error)
	ModifyProfile(user *User, ctxUser web.User) (err error)
}

func NewUser(d dao.Interface, eb EventBiz) UserBiz {
	return &userBiz{d: d, eb: eb}
}

type userBiz struct {
	d  dao.Interface
	eb EventBiz
}

func (b *userBiz) Search(name, loginName, filter string, pageIndex, pageSize int) (users []*User, total int, err error) {
	var (
		list []*model.User
		args = &model.UserSearchArgs{
			Name:      name,
			LoginName: loginName,
			Status:    -1,
			PageIndex: pageIndex,
			PageSize:  pageSize,
		}
	)

	switch filter {
	case "admins":
		args.Admin = true
	case "active":
		args.Status = UserStatusActive
	case "blocked":
		args.Status = UserStatusBlocked
	}

	list, total, err = b.d.UserSearch(context.TODO(), args)
	if err == nil {
		for _, u := range list {
			users = append(users, newUser(u))
		}
	}
	return
}

func (b *userBiz) FindByID(id string) (user *User, err error) {
	var u *model.User
	u, err = b.d.UserGet(context.TODO(), id)
	if u != nil {
		user = newUser(u)
	}
	return
}

func (b *userBiz) FindByName(loginName string) (user *User, err error) {
	var u *model.User
	u, err = b.d.UserGetByName(context.TODO(), loginName)
	if u != nil {
		user = newUser(u)
	}
	return
}

func (b *userBiz) FindPrivacy(loginName string) (privacy *UserPrivacy, err error) {
	var u *model.User
	u, err = b.d.UserGetByName(context.TODO(), loginName)
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

func (b *userBiz) Create(u *User, ctxUser web.User) (id string, err error) {
	user := &model.User{
		ID:        createId(),
		Name:      u.Name,
		LoginName: u.LoginName,
		Email:     u.Email,
		Admin:     u.Admin,
		Type:      u.Type,
		Status:    UserStatusActive,
		Roles:     u.Roles,
		CreatedAt: time.Now(),
		CreatedBy: model.Operator{ID: ctxUser.ID(), Name: ctxUser.Name()},
	}
	user.UpdatedAt = user.CreatedAt
	user.UpdatedBy = user.CreatedBy
	if user.Type == UserTypeInternal {
		user.Password, user.Salt, err = passwd.Generate(u.Password)
		if err != nil {
			return
		}
	}

	if err = b.d.UserCreate(context.TODO(), user); err == nil && ctxUser != nil {
		b.eb.CreateUser(EventActionCreate, user.LoginName, user.Name, ctxUser)
	}
	id = user.ID
	return
}

func (b *userBiz) Update(u *User, ctxUser web.User) (err error) {
	user := &model.User{
		ID:        u.ID,
		Name:      u.Name,
		LoginName: u.LoginName,
		Email:     u.Email,
		Admin:     u.Admin,
		Type:      u.Type,
		Roles:     u.Roles,
		UpdatedAt: time.Now(),
	}
	user.UpdatedBy.ID = ctxUser.ID()
	user.UpdatedBy.Name = ctxUser.Name()
	if err = b.d.UserUpdate(context.TODO(), user); err == nil {
		b.eb.CreateUser(EventActionUpdate, u.LoginName, u.Name, ctxUser)
	}
	return
}

func (b *userBiz) SetStatus(id string, status int32, user web.User) (err error) {
	u := &model.User{
		ID:        id,
		Status:    status,
		UpdatedAt: time.Now(),
		UpdatedBy: model.Operator{ID: user.ID(), Name: user.Name()},
	}
	return b.d.UserUpdateStatus(context.TODO(), u)
}

func (b *userBiz) Delete(id, name string, user web.User) (err error) {
	err = b.d.UserDelete(context.TODO(), id)
	if err == nil {
		b.eb.CreateUser(EventActionDelete, id, name, user)
	}
	return
}

func (b *userBiz) ModifyPassword(oldPwd, newPwd string, user web.User) (err error) {
	var u *model.User
	u, err = b.d.UserGet(context.TODO(), user.ID())
	if err != nil {
		return err
	} else if u == nil {
		return errors.Format("user not found: %s", user.ID())
	}

	if !passwd.Validate(oldPwd, u.Password, u.Salt) {
		return errors.Coded(misc.ErrOldPasswordIncorrect, "current password is incorrect")
	}

	u.Password, u.Salt, err = passwd.Generate(newPwd)
	if err != nil {
		return
	}

	u.UpdatedAt = time.Now()
	u.UpdatedBy.ID = user.ID()
	u.UpdatedBy.Name = user.Name()
	err = b.d.UserUpdatePassword(context.TODO(), u)
	return
}

func (b *userBiz) ModifyProfile(u *User, user web.User) (err error) {
	return b.d.UserUpdateProfile(context.TODO(), &model.User{
		ID:        user.ID(),
		Name:      u.Name,
		LoginName: u.LoginName,
		Email:     u.Email,
		UpdatedAt: time.Now(),
		UpdatedBy: model.Operator{ID: user.ID(), Name: user.Name()},
	})
}

func (b *userBiz) Count() (count int, err error) {
	return b.d.UserCount(context.TODO())
}

func (b *userBiz) UpdateSession(id string) (token string, err error) {
	session := &model.Session{
		UserID:    id,
		Token:     guid.New().String(),
		UpdatedAt: time.Now(),
	}
	session.Expires = session.UpdatedAt.Add(time.Hour * 24)
	err = b.d.SessionUpdate(context.TODO(), session)
	if err == nil {
		token = session.Token
	}
	return
}

func (b *userBiz) GetSession(token string) (session *model.Session, err error) {
	return b.d.SessionGet(context.TODO(), token)
}

type User struct {
	ID        string         `json:"id"`
	Name      string         `json:"name" valid:"required"`
	LoginName string         `json:"loginName" valid:"required"`
	Password  string         `json:"password,omitempty"`
	Email     string         `json:"email" valid:"required"`
	Admin     bool           `json:"admin"`
	Type      string         `json:"type"`
	Status    int32          `json:"status,omitempty"`
	Roles     []string       `json:"roles,omitempty"`
	CreatedAt string         `json:"createdAt,omitempty"`
	UpdatedAt string         `json:"updatedAt,omitempty"`
	CreatedBy model.Operator `json:"createdBy"`
	UpdatedBy model.Operator `json:"updatedBy"`
}

type UserPrivacy struct {
	ID       string
	Name     string
	Password string
	Salt     string
	Type     string
	Status   int32
}

func newUser(u *model.User) *User {
	user := &User{
		CreatedAt: formatTime(u.CreatedAt),
		UpdatedAt: formatTime(u.UpdatedAt),
	}
	_ = copier.CopyWithOption(user, u, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	return user
}

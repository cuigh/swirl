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

type User struct {
	ID        string   `json:"id"`
	Name      string   `json:"name" valid:"required"`
	LoginName string   `json:"loginName" valid:"required"`
	Password  string   `json:"password,omitempty"`
	Email     string   `json:"email" valid:"required"`
	Admin     bool     `json:"admin"`
	Type      string   `json:"type"`
	Status    int32    `json:"status,omitempty"`
	Roles     []string `json:"roles,omitempty"`
	CreatedAt string   `bson:"created_at" json:"createdAt,omitempty"`
	UpdatedAt string   `bson:"updated_at" json:"updatedAt,omitempty"`
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

type UserBiz interface {
	Search(name, loginName, filter string, pageIndex, pageSize int) (users []*User, total int, err error)
	Create(user *User, ctxUser web.User) (id string, err error)
	Update(user *User, ctxUser web.User) (err error)
	FindByID(id string) (user *User, err error)
	FindByName(loginName string) (user *User, err error)
	FindPrivacy(loginName string) (privacy *UserPrivacy, err error)
	Count() (count int, err error)
	Delete(id, name string, user web.User) (err error)
	SetStatus(id string, status int32) (err error)
	ModifyPassword(id, pwd, salt string) (err error)
	ModifyProfile(user *User) (err error)
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

	list, total, err = b.d.UserList(context.TODO(), args)
	if err == nil {
		for _, u := range list {
			users = append(users, newUser(u))
		}
	}
	return
}

func (b *userBiz) FindByID(id string) (user *User, err error) {
	var u *model.User
	u, err = b.d.UserGetByID(context.TODO(), id)
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
	}
	user.UpdatedAt = user.CreatedAt
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
	if err = b.d.UserUpdate(context.TODO(), user); err == nil {
		b.eb.CreateUser(EventActionUpdate, u.LoginName, u.Name, ctxUser)
	}
	return
}

func (b *userBiz) SetStatus(id string, status int32) (err error) {
	return b.d.UserSetStatus(context.TODO(), id, status)
}

func (b *userBiz) Delete(id, name string, user web.User) (err error) {
	err = b.d.UserDelete(context.TODO(), id)
	if err == nil {
		b.eb.CreateUser(EventActionDelete, id, name, user)
	}
	return
}

func (b *userBiz) ModifyPassword(id, oldPwd, newPwd string) (err error) {
	var (
		user      *model.User
		pwd, salt string
	)

	user, err = b.d.UserGetByID(context.TODO(), id)
	if err != nil {
		return err
	} else if user == nil {
		return errors.Format("user not found: %s", id)
	}

	if !passwd.Validate(oldPwd, user.Password, user.Salt) {
		return errors.Coded(misc.ErrOldPasswordIncorrect, "current password is incorrect")
	}

	pwd, salt, err = passwd.Generate(newPwd)
	if err != nil {
		return
	}

	err = b.d.UserModifyPassword(context.TODO(), id, pwd, salt)
	return
}

func (b *userBiz) ModifyProfile(user *User) (err error) {
	return b.d.UserModifyProfile(context.TODO(), &model.User{
		ID:        user.ID,
		Name:      user.Name,
		LoginName: user.LoginName,
		Email:     user.Email,
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

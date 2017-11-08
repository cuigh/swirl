package biz

import (
	"fmt"
	"time"

	"github.com/cuigh/auxo/data/guid"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/security/password"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
	"github.com/go-ldap/ldap"
)

var ErrIncorrectAuth = errors.New("login name or password is incorrect")

var User = &userBiz{}

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
	user.ID = guid.New()
	user.Status = model.UserStatusActive
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	if user.Type == model.UserTypeInternal {
		user.Password, user.Salt, err = password.Get(user.Password)
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

func (b *userBiz) UpdatePassword(id, old_pwd, new_pwd string) (err error) {
	do(func(d dao.Interface) {
		var (
			user      *model.User
			pwd, salt string
		)

		user, err = d.UserGetByID(id)
		if err != nil {
			return
		}

		if !password.Validate(user.Password, old_pwd, user.Salt) {
			err = errors.New("Current password is incorrect")
			return
		}

		pwd, salt, err = password.Get(new_pwd)
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

func (b *userBiz) Login(name, pwd string) (token string, err error) {
	do(func(d dao.Interface) {
		var (
			user *model.User
		)

		user, err = d.UserGetByName(name)
		if err != nil {
			return
		}

		if user == nil {
			user = &model.User{
				Type:      model.UserTypeLDAP,
				LoginName: name,
			}
			err = b.loginLDAP(user, pwd)
		} else {
			if user.Status == model.UserStatusBlocked {
				err = fmt.Errorf("user %s is blocked", name)
				return
			}

			if user.Type == model.UserTypeInternal {
				err = b.loginInternal(user, pwd)
			} else {
				err = b.loginLDAP(user, pwd)
			}
		}

		if err != nil {
			return
		}

		session := &model.Session{
			UserID:    user.ID,
			Token:     guid.New(),
			UpdatedAt: time.Now(),
		}
		session.Expires = session.UpdatedAt.Add(time.Hour * 24)
		err = d.SessionUpdate(session)
		if err != nil {
			return
		}
		token = session.Token

		// create event
		Event.CreateAuthentication(model.EventActionLogin, user.ID, user.LoginName, user.Name)
	})
	return
}

func (b *userBiz) loginInternal(user *model.User, pwd string) error {
	if !password.Validate(user.Password, pwd, user.Salt) {
		return ErrIncorrectAuth
	}
	return nil
}

func (b *userBiz) loginLDAP(user *model.User, pwd string) error {
	setting, err := Setting.Get()
	if err != nil {
		return err
	}

	if !setting.LDAP.Enabled {
		return ErrIncorrectAuth
	}

	l, err := ldap.Dial("tcp", setting.LDAP.Address)
	if err != nil {
		return err
	}
	defer l.Close()

	// bind
	err = l.Bind(user.LoginName, pwd)
	if err != nil {
		log.Get("user").Error("Login by LDAP failed: ", err)
		return ErrIncorrectAuth
	}

	// Stop here for an exist user because we only need validate password.
	if user.ID != "" {
		return nil
	}

	// If user wasn't exist, we need create it
	req := ldap.NewSearchRequest(
		setting.LDAP.BaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(%s=%s))", setting.LDAP.LoginAttr, user.LoginName),
		[]string{"dn", setting.LDAP.EmailAttr, setting.LDAP.LoginAttr, setting.LDAP.NameAttr},
		nil,
	)
	searchResult, err := l.Search(req)
	if err != nil {
		return err
	}
	if len(searchResult.Entries) == 0 {
		return ErrIncorrectAuth
	}

	entry := searchResult.Entries[0]
	user.Email = entry.GetAttributeValue(setting.LDAP.EmailAttr)
	user.Name = entry.GetAttributeValue(setting.LDAP.NameAttr)
	if user.ID == "" {
		return b.Create(user, nil)
	}
	return nil
}

// Identify authenticate user
func (b *userBiz) Identify(token string) (user web.User) {
	do(func(d dao.Interface) {
		var (
			roles []*model.Role
			role  *model.Role
		)

		session, err := d.SessionGet(token)
		if err != nil {
			log.Get("user").Errorf("Load session failed: %v", err)
			return
		}
		if session == nil || session.Expires.Before(time.Now()) {
			return
		}

		u, err := d.UserGetByID(session.UserID)
		if err != nil {
			log.Get("user").Errorf("Load user failed: %v", err)
			return
		}
		if u == nil {
			return
		}

		if len(u.Roles) > 0 {
			roles = make([]*model.Role, len(u.Roles))
			for i, id := range u.Roles {
				role, err = d.RoleGet(id)
				if err != nil {
					return
				} else if role != nil {
					roles[i] = role
				}
			}
		}
		user = model.NewAuthUser(u, roles)
	})
	return
}

// Authorize check permission of user
func (b *userBiz) Authorize(user web.User, h web.HandlerInfo) bool {
	if au, ok := user.(*model.AuthUser); ok {
		return au.IsAllowed(h.Name())
	}
	return false
}

//
//func (b *userBiz) Find(key string) string {
//	b.locker.Lock()
//	defer b.locker.Unlock()
//
//	return b.tickets[key]
//}
//
//func (b *userBiz) setTicket(name string) (key string) {
//	b.locker.Lock()
//	defer b.locker.Unlock()
//
//	key = guid.New()
//	b.tickets[key] = name
//	return
//}

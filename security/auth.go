package security

import (
	"time"

	"github.com/cuigh/auxo/cache"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/security"
	"github.com/cuigh/auxo/security/certify"
	"github.com/cuigh/auxo/security/certify/ldap"
	"github.com/cuigh/auxo/security/passwd"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

const pkgName = "swirl.security"

// Validator is the authenticator function of swirl.
func Validator(setting *model.Setting) func(name, pwd string) (ticket string, err error) {
	ldapRealm := createLDAPRealm(setting)
	return func(name, pwd string) (ticket string, err error) {
		var (
			su security.User
			mu *model.User
		)

		// try find user first
		mu, err = biz.User.GetByName(name)
		if err != nil {
			return
		}

		if mu != nil && mu.Type == model.UserTypeInternal { // internal user
			if !passwd.Validate(pwd, mu.Password, mu.Salt) {
				err = certify.ErrInvalidToken
			}
		} else if ldapRealm != nil { // user not exist or ldap user
			su, err = ldapRealm.Login(certify.NewSimpleToken(name, pwd))
			if err == nil && mu == nil { // create user if not exists
				lu := su.(*ldap.User)
				mu = &model.User{
					Type:      model.UserTypeLDAP,
					LoginName: lu.LoginName(),
					Name:      lu.Name(),
					Email:     lu.Email(),
				}
				err = biz.User.Create(mu, nil)
			}
		} else {
			err = certify.ErrInvalidToken
		}

		// replace user session, one session is allowed per user.
		if err == nil {
			ticket, err = biz.User.UpdateSession(mu.ID)
		}

		// create event
		if err == nil {
			biz.Event.CreateAuthentication(model.EventActionLogin, mu.ID, mu.LoginName, mu.Name)
		}
		return
	}
}

func createLDAPRealm(setting *model.Setting) certify.Realm {
	if setting.LDAP.Enabled {
		opts := []ldap.Option{
			ldap.NameAttr(setting.LDAP.NameAttr),
			ldap.EmailAttr(setting.LDAP.EmailAttr),
			ldap.UserFilter(setting.LDAP.UserFilter),
			ldap.Security(ldap.SecurityPolicy(setting.LDAP.Security)),
		}
		if setting.LDAP.Authentication == 1 {
			opts = append(opts, ldap.Binding(setting.LDAP.BindDN, setting.LDAP.BindPassword))
		}
		return ldap.New(setting.LDAP.Address, setting.LDAP.BaseDN, setting.LDAP.UserDN, opts...)
	}
	return nil
}

// Identifier is used to identity user.
func Identifier(token string) (user web.User) {
	const cacheKey = "auth_user"

	session, err := biz.User.GetSession(token)
	if err != nil {
		log.Get(pkgName).Errorf("Load session failed: %v", err)
		return
	}
	if session == nil || session.Expires.Before(time.Now()) {
		return
	}

	// try find from cache first
	value := cache.Get(cacheKey, session.UserID)
	if !value.IsNil() {
		user = &model.AuthUser{}
		if err = value.Scan(user); err == nil {
			return
		}
		log.Get(pkgName).Warnf("Load auth user from cache failed: %v", err)
	}

	u, err := biz.User.GetByID(session.UserID)
	if err != nil {
		log.Get(pkgName).Errorf("Load user failed: %v", err)
		return
	}
	if u == nil {
		return
	}

	var roles []*model.Role
	if roles, err = getRoles(u); err == nil {
		user = model.NewAuthUser(u, roles)
		cache.Set(user, cacheKey, session.UserID)
	}
	return
}

func getRoles(u *model.User) (roles []*model.Role, err error) {
	if len(u.Roles) > 0 {
		roles = make([]*model.Role, len(u.Roles))
		for i, id := range u.Roles {
			roles[i], err = biz.Role.Get(id)
			if err != nil {
				return nil, err
			}
		}
	}
	return
}

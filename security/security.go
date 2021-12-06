package security

import (
	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/security"
	"github.com/cuigh/auxo/security/certify"
	"github.com/cuigh/auxo/security/certify/ldap"
	"github.com/cuigh/auxo/security/passwd"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/misc"
)

type Authenticator struct {
	ub     biz.UserBiz
	realms []RealmFunc
}

func NewAuthenticator(s *misc.Setting, ub biz.UserBiz) *Authenticator {
	return &Authenticator{
		ub:     ub,
		realms: []RealmFunc{internalRealm(), ldapRealm(s, ub)},
	}
}

func (a *Authenticator) Login(loginName, password string) (user security.User, err error) {
	privacy, err := a.ub.FindPrivacy(loginName)
	if err != nil {
		return nil, err
	}

	if privacy != nil && privacy.Status == biz.UserStatusBlocked {
		return nil, misc.Error(misc.ErrAccountDisabled, certify.ErrAccountDisabled)
	}

	for _, login := range a.realms {
		user, err = login(privacy, loginName, password)
		if user != nil && err == nil {
			return
		}
	}

	if err == nil {
		err = misc.Error(misc.ErrInvalidToken, certify.ErrInvalidToken)
	}
	return
}

type RealmFunc func(u *biz.UserPrivacy, loginName, password string) (security.User, error)

func internalRealm() RealmFunc {
	return func(u *biz.UserPrivacy, loginName, password string) (security.User, error) {
		if u == nil || u.Type != biz.UserTypeInternal {
			return nil, nil
		}

		if passwd.Validate(password, u.Password, u.Salt) {
			return security.NewUser(u.ID, u.Name), nil
		}
		return nil, misc.Error(misc.ErrInvalidToken, certify.ErrInvalidToken)
	}
}

func ldapRealm(s *misc.Setting, ub biz.UserBiz) RealmFunc {
	const authBind = "bind"
	var r certify.Realm

	if s.LDAP.Enabled {
		opts := []ldap.Option{
			ldap.NameAttr(s.LDAP.NameAttr),
			ldap.EmailAttr(s.LDAP.EmailAttr),
			ldap.UserFilter(s.LDAP.UserFilter),
			ldap.Security(ldap.SecurityPolicy(s.LDAP.Security)),
		}
		if s.LDAP.Authentication == authBind {
			opts = append(opts, ldap.Binding(s.LDAP.BindDN, s.LDAP.BindPassword))
		}
		r = ldap.New(s.LDAP.Address, s.LDAP.BaseDN, s.LDAP.UserDN, opts...)
	}

	return func(u *biz.UserPrivacy, loginName, password string) (security.User, error) {
		if r == nil || (u != nil && u.Type != biz.UserTypeLDAP) {
			return nil, nil
		}

		user, err := r.Login(certify.NewSimpleToken(loginName, password))
		if err != nil {
			return nil, err
		}

		var (
			id string
			lu = user.(*ldap.User)
		)
		if u == nil {
			id, err = ub.Create(&biz.User{
				Type:      biz.UserTypeLDAP,
				LoginName: loginName,
				Name:      lu.Name(),
				Email:     lu.Email(),
			}, nil)
			if err != nil {
				return nil, err
			}
			lu.SetID(id)
		} else {
			lu.SetID(u.ID)
		}
		return user, nil
	}
}

func init() {
	container.Put(NewAuthenticator)
	container.Put(NewIdentifier, container.Name("identifier"))
	container.Put(NewAuthorizer, container.Name("authorizer"))
}

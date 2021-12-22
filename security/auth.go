package security

import (
	"strings"
	"time"

	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/security"
	"github.com/cuigh/auxo/security/certify"
	"github.com/cuigh/auxo/security/certify/ldap"
	"github.com/cuigh/auxo/security/passwd"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Identifier identifies the user.
type Identifier struct {
	ub         biz.UserBiz
	rb         biz.RoleBiz
	sb         biz.SessionBiz
	realms     []RealmFunc
	extractors []TokenExtractor
	logger     log.Logger
}

func NewIdentifier(s *misc.Setting, ub biz.UserBiz, rb biz.RoleBiz, sb biz.SessionBiz) *Identifier {
	return &Identifier{
		ub:         ub,
		rb:         rb,
		sb:         sb,
		realms:     []RealmFunc{internalRealm(), ldapRealm(s, ub)},
		extractors: []TokenExtractor{headerExtractor, queryExtractor},
		logger:     log.Get(PkgName),
	}
}

func (c *Identifier) Apply(next web.HandlerFunc) web.HandlerFunc {
	return func(ctx web.Context) error {
		token := c.extractToken(ctx)
		if token != "" {
			user := c.identifyUser(token)
			ctx.SetUser(user)
		}
		return next(ctx)
	}
}

func (c *Identifier) Identify(loginName, password string) (identify Identity, err error) {
	var (
		u security.User
		s *model.Session
	)

	u, err = c.signIn(loginName, password)
	if err != nil {
		return
	}

	s, err = c.createSession(u)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		id:    u.ID(),
		name:  u.Name(),
		token: s.ID,
		perms: s.Perms,
	}, nil
}

func (c *Identifier) signIn(loginName, password string) (user security.User, err error) {
	privacy, err := c.ub.FindPrivacy(loginName)
	if err != nil {
		return nil, err
	}

	if privacy != nil && privacy.Status == biz.UserStatusBlocked {
		return nil, misc.Error(misc.ErrAccountDisabled, certify.ErrAccountDisabled)
	}

	for _, login := range c.realms {
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

func (c *Identifier) extractToken(ctx web.Context) (token string) {
	for _, e := range c.extractors {
		if token = e(ctx); token != "" {
			return
		}
	}
	return
}

func (c *Identifier) identifyUser(token string) web.User {
	session, err := c.sb.Find(token)
	if err != nil {
		c.logger.Error("failed to find session: ", err)
		return nil
	} else if session == nil || session.Expiry.Before(time.Now()) {
		return nil
	}

	if session.Dirty {
		if err = c.updateSession(session); err != nil {
			c.logger.Error("failed to refresh session: ", err)
			return nil
		}
	} else if time.Now().Add(time.Minute * 5).After(session.Expiry) {
		c.renewSession(session)
	}

	return c.createUser(session)
}

func (c *Identifier) createUser(s *model.Session) web.User {
	return &User{
		token: s.ID,
		id:    s.UserID,
		name:  s.Username,
		admin: s.Admin,
		perm:  PermMap(s.Perm),
	}
}

func (c *Identifier) createSession(user security.User) (s *model.Session, err error) {
	s = &model.Session{
		ID:       primitive.NewObjectID().Hex(),
		UserID:   user.ID(),
		Username: user.Name(),
		Expiry:   time.Now().Add(misc.Options.TokenExpiry),
	}
	s.MaxExpiry = s.Expiry.Add(24 * time.Hour)
	if err = c.fillSession(s); err == nil {
		err = c.sb.Create(s)
	}
	return
}

func (c *Identifier) updateSession(s *model.Session) (err error) {
	if err = c.fillSession(s); err == nil {
		err = c.sb.Update(s)
	}
	return
}

func (c *Identifier) fillSession(s *model.Session) (err error) {
	u, err := c.ub.FindByID(s.UserID)
	if err != nil {
		return err
	} else if u == nil {
		return errors.New("user not found")
	}

	if u.Admin {
		s.Perms = []string{"*"}
	} else {
		s.Perms, err = c.rb.GetPerms(u.Roles)
		if err != nil {
			return err
		}
		s.Perm = uint64(NewPermMap(s.Perms))
	}

	s.Username = u.Name
	s.Admin = u.Admin
	s.Roles = u.Roles
	return nil
}

func (c *Identifier) renewSession(s *model.Session) {
	expiry := time.Now().Add(misc.Options.TokenExpiry)
	if expiry.After(s.MaxExpiry) {
		expiry = s.MaxExpiry
	}
	err := c.sb.UpdateExpiry(s.ID, expiry)
	if err != nil {
		c.logger.Errorf("failed to renew token '%s': %s", s.ID, err)
	}
}

type TokenExtractor func(ctx web.Context) string

func headerExtractor(ctx web.Context) (token string) {
	const prefix = "Bearer "
	if value := ctx.Header(web.HeaderAuthorization); strings.HasPrefix(value, prefix) {
		token = value[len(prefix):]
	}
	return
}

func queryExtractor(ctx web.Context) (token string) {
	return ctx.Query("token")
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
			id, err = ub.Create(&model.User{
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

type Identity interface {
	ID() string
	Name() string
	Anonymous() bool
	Token() string
	Perms() []string
}

type UserInfo struct {
	id    string
	name  string
	token string
	perms []string
}

func (u *UserInfo) ID() string {
	return u.id
}

func (u *UserInfo) Name() string {
	return u.name
}

func (u *UserInfo) Anonymous() bool {
	return u.id == ""
}

func (u *UserInfo) Token() string {
	return u.token
}

func (u *UserInfo) Perms() []string {
	return u.perms
}

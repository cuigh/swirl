package security

import (
	"net/http"
	"time"

	"github.com/cuigh/auxo/cache"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
)

type Authorizer struct {
	ub     biz.UserBiz
	perms  *cache.Value
	logger log.Logger
}

func NewAuthorizer(ub biz.UserBiz, rb biz.RoleBiz) web.Filter {
	v := cache.Value{
		TTL:  5 * time.Minute,
		Load: func() (interface{}, error) { return loadPerms(rb) },
	}
	return &Authorizer{
		ub:     ub,
		perms:  &v,
		logger: log.Get("security"),
	}
}

// Apply implements `web.Filter` interface.
func (a *Authorizer) Apply(next web.HandlerFunc) web.HandlerFunc {
	return func(ctx web.Context) error {
		auth := ctx.Handler().Authorize()

		// allow anonymous
		if auth == "" || auth == web.AuthAnonymous {
			return next(ctx)
		}

		user := ctx.User()
		if user == nil || user.Anonymous() {
			return web.NewError(http.StatusUnauthorized, "You are not logged in")
		}

		if auth != web.AuthAuthenticated && !a.check(user, auth) {
			return web.NewError(http.StatusForbidden, "You do not have access to this resource")
		}
		return next(ctx)
	}
}

func (a *Authorizer) check(user web.User, auth string) bool {
	u, err := a.ub.FindByID(user.ID())
	if err != nil {
		a.logger.Errorf("failed to query user '%s': %s", user.ID(), err)
		return false
	}

	if u == nil || u.Status == biz.UserStatusBlocked {
		return false
	} else if u.Admin {
		return true
	} else if auth == web.AuthAdministrator {
		return u.Admin
	}

	v, err := a.perms.Get(true)
	if err != nil {
		a.logger.Error("failed to load role perms: ", err)
		return false
	}

	perms := v.(map[string]PermSet)
	for _, r := range u.Roles {
		if set, ok := perms[r]; ok {
			if set.Contains(auth) {
				return true
			}
		}
	}
	return false
}

func loadPerms(rb biz.RoleBiz) (interface{}, error) {
	roles, err := rb.Search("")
	if err != nil {
		return nil, err
	}

	perms := make(map[string]PermSet)
	for _, role := range roles {
		set := make(PermSet)
		for _, p := range role.Perms {
			set[p] = struct{}{}
		}
		perms[role.ID] = set
	}
	return perms, nil
}

type PermSet map[string]struct{}

func (s PermSet) Contains(perm string) (ok bool) {
	_, ok = s[perm]
	return
}

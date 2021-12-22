package security

import (
	"net/http"
	"strings"

	"github.com/cuigh/auxo/net/web"
)

const ActionBits = 24

// Resources holds all resources requiring authorization. Up to 40 resources are supported.
// WARN: DO NOT CHANGE VALUES!!!
var Resources = map[string]uint64{
	"registry":  1,
	"node":      1 << 1,
	"network":   1 << 2,
	"service":   1 << 3,
	"task":      1 << 4,
	"stack":     1 << 5,
	"config":    1 << 6,
	"secret":    1 << 7,
	"image":     1 << 8,
	"container": 1 << 9,
	"volume":    1 << 10,
	"user":      1 << 11,
	"role":      1 << 12,
	"chart":     1 << 13,
	"dashboard": 1 << 14,
	"event":     1 << 15,
	"setting":   1 << 16,
}

// Actions holds all actions requiring authorization. Up to 24 actions are supported.
// WARN: DO NOT CHANGE VALUES!!!
var Actions = map[string]uint64{
	"view":       1,
	"edit":       1 << 1,
	"delete":     1 << 2,
	"disconnect": 1 << 3,
	"restart":    1 << 4,
	"rollback":   1 << 5,
	"logs":       1 << 6,
	"deploy":     1 << 7,
	"shutdown":   1 << 8,
	"execute":    1 << 9,
}

var Perms = map[string][]string{
	"registry":  {"view", "edit", "delete"},
	"node":      {"view", "edit", "delete"},
	"network":   {"view", "edit", "delete", "disconnect"},
	"service":   {"view", "edit", "delete", "restart", "rollback", "logs"},
	"task":      {"view", "logs"},
	"stack":     {"view", "edit", "delete", "deploy", "shutdown"},
	"config":    {"view", "edit", "delete"},
	"secret":    {"view", "edit", "delete"},
	"image":     {"view", "delete"},
	"container": {"view", "delete", "logs", "execute"},
	"volume":    {"view", "edit", "delete"},
	"user":      {"view", "edit", "delete"},
	"role":      {"view", "edit", "delete"},
	"chart":     {"view", "edit", "delete"},
	"dashboard": {"edit"},
	"event":     {"view"},
	"setting":   {"view", "edit"},
}

type Authorizer struct {
}

func NewAuthorizer() Authorizer {
	return Authorizer{}
}

// Apply implements `web.Filter` interface.
func (p Authorizer) Apply(next web.HandlerFunc) web.HandlerFunc {
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

		if auth != web.AuthAuthenticated && !p.check(user, auth) {
			return web.NewError(http.StatusForbidden, "You do not have access to this resource")
		}
		return next(ctx)
	}
}

func (p Authorizer) check(user web.User, auth string) bool {
	u := user.(*User)
	if u.admin {
		return true
	}
	return u.perm.Contains(auth)
}

type PermMap uint64

func NewPermMap(perms []string) PermMap {
	var p uint64
	for _, perm := range perms {
		pair := strings.SplitN(perm, ".", 2)
		if len(pair) == 2 {
			r, a := Resources[pair[0]], Actions[pair[1]]
			p |= r << ActionBits
			p |= a
		}
	}
	return PermMap(p)
}

func (p PermMap) Contains(perm string) bool {
	pair := strings.SplitN(perm, ".", 2)
	if len(pair) == 2 {
		r, a := Resources[pair[0]], Actions[pair[1]]
		return uint64(p)&(r<<ActionBits) > 0 && uint64(p)&a > 0
	}
	return false
}

type User struct {
	token string
	id    string
	name  string
	admin bool
	perm  PermMap
}

func (u *User) Token() string {
	return u.token
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Anonymous() bool {
	return u.id == ""
}

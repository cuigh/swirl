package model

import (
	"time"
)

type UserType string

const (
	// UserTypeInternal is internal user of swirl
	UserTypeInternal UserType = "internal"
	// UserTypeLDAP is external user of LDAP
	UserTypeLDAP UserType = "ldap"
)

var Placeholder = struct{}{}

type UserStatus int32

const (
	// UserStatusBlocked is the status which user is blocked
	UserStatusBlocked UserStatus = 0
	// UserStatusActive is the normal status
	UserStatusActive UserStatus = 1
)

type Role struct {
	ID          string    `bson:"_id" json:"id,omitempty"`
	Name        string    `bson:"name" json:"name,omitempty" valid:"required"`
	Description string    `bson:"desc" json:"desc,omitempty"`
	Perms       []string  `bson:"perms" json:"perms,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at,omitempty"`
}

type User struct {
	ID        string     `bson:"_id" json:"id,omitempty"`
	Name      string     `bson:"name" json:"name,omitempty" valid:"required"`
	LoginName string     `bson:"login_name" json:"login_name,omitempty" valid:"required"`
	Password  string     `bson:"password" json:"password,omitempty"`
	Salt      string     `bson:"salt" json:"salt,omitempty"`
	Email     string     `bson:"email" json:"email,omitempty" valid:"required"`
	Admin     bool       `bson:"admin" json:"admin,omitempty"`
	Type      UserType   `bson:"type" json:"type,omitempty"`
	Status    UserStatus `bson:"status" json:"status,omitempty"`
	Roles     []string   `bson:"roles" json:"roles,omitempty"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at,omitempty"`
}

type UserListArgs struct {
	// admins, active, blocked
	Filter    string `bind:"filter"`
	Query     string `bind:"query"`
	PageIndex int    `bind:"page"`
	PageSize  int    `bind:"size"`
}

type Session struct {
	UserID    string    `bson:"_id" json:"id,omitempty"`
	Token     string    `bson:"token" json:"token,omitempty"`
	Expires   time.Time `bson:"expires" json:"expires,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
}

type AuthUser struct {
	user  *User
	perms map[string]struct{}
}

func NewAuthUser(user *User, roles []*Role) *AuthUser {
	if user == nil {
		panic(111)
	}
	u := &AuthUser{
		user:  user,
		perms: make(map[string]struct{}),
	}
	for _, role := range roles {
		for _, perm := range role.Perms {
			u.perms[perm] = Placeholder
		}
	}
	return u
}

func (u *AuthUser) ID() string {
	return u.user.ID
}

func (u *AuthUser) Name() string {
	return u.user.Name
}

func (u *AuthUser) Anonymous() bool {
	return u.user.ID == ""
}

func (u *AuthUser) Admin() bool {
	return u.user.Admin
}

func (u *AuthUser) IsAllowed(perm string) bool {
	if u.user.Admin {
		return true
	}

	_, ok := u.perms[perm]
	return ok
}

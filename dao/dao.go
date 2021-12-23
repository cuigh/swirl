package dao

import (
	"context"
	"time"

	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/misc"
)

var builders = make(map[string]Builder)

// Builder creates an Interface instance.
type Builder func(addr string) (Interface, error)

func Register(name string, builder Builder) {
	builders[name] = builder
}

// Interface is the interface that wraps all dao methods.
type Interface interface {
	Init() error

	RoleGet(ctx context.Context, id string) (*Role, error)
	RoleSearch(ctx context.Context, name string) (roles []*Role, err error)
	RoleCreate(ctx context.Context, role *Role) error
	RoleUpdate(ctx context.Context, role *Role) error
	RoleDelete(ctx context.Context, id string) error

	UserGet(ctx context.Context, id string) (*User, error)
	UserGetByName(ctx context.Context, loginName string) (*User, error)
	UserSearch(ctx context.Context, args *UserSearchArgs) (users []*User, count int, err error)
	UserCount(ctx context.Context) (int, error)
	UserCreate(ctx context.Context, user *User) error
	UserUpdate(ctx context.Context, user *User) error
	UserUpdateStatus(ctx context.Context, user *User) error
	UserUpdateProfile(ctx context.Context, user *User) error
	UserUpdatePassword(ctx context.Context, user *User) error
	UserDelete(ctx context.Context, id string) error

	SessionGet(ctx context.Context, id string) (*Session, error)
	SessionCreate(ctx context.Context, session *Session) error
	SessionUpdate(ctx context.Context, session *Session) error
	SessionUpdateExpiry(ctx context.Context, id string, expiry time.Time) (err error)
	SessionUpdateDirty(ctx context.Context, userID string, roleID string) (err error)

	RegistryGet(ctx context.Context, id string) (*Registry, error)
	RegistryGetByURL(ctx context.Context, url string) (registry *Registry, err error)
	RegistryGetAll(ctx context.Context) (registries []*Registry, err error)
	RegistryCreate(ctx context.Context, registry *Registry) error
	RegistryUpdate(ctx context.Context, registry *Registry) error
	RegistryDelete(ctx context.Context, id string) error

	StackGet(ctx context.Context, name string) (*Stack, error)
	StackGetAll(ctx context.Context) (stacks []*Stack, err error)
	StackCreate(ctx context.Context, stack *Stack) error
	StackUpdate(ctx context.Context, stack *Stack) error
	StackDelete(ctx context.Context, name string) error

	EventSearch(ctx context.Context, args *EventSearchArgs) (events []*Event, count int, err error)
	EventCreate(ctx context.Context, event *Event) error

	SettingGet(ctx context.Context, id string) (*Setting, error)
	SettingGetAll(ctx context.Context) (settings []*Setting, err error)
	SettingUpdate(ctx context.Context, setting *Setting) error

	ChartGet(ctx context.Context, id string) (*Chart, error)
	ChartGetBatch(ctx context.Context, ids ...string) ([]*Chart, error)
	ChartSearch(ctx context.Context, args *ChartSearchArgs) (charts []*Chart, count int, err error)
	ChartCreate(ctx context.Context, chart *Chart) error
	ChartUpdate(ctx context.Context, chart *Chart) error
	ChartDelete(ctx context.Context, id string) error

	DashboardGet(ctx context.Context, name, key string) (dashboard *Dashboard, err error)
	DashboardUpdate(ctx context.Context, dashboard *Dashboard) error
}

func newInterface() (i Interface) {
	var err error

	if b, ok := builders[misc.Options.DBType]; ok {
		i, err = b(misc.Options.DBAddress)
	} else {
		err = errors.New("unknown database type: " + misc.Options.DBType)
	}

	if err != nil {
		panic(err)
	}
	return
}

func init() {
	container.Put(newInterface)
}

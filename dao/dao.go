package dao

import (
	"context"
	"time"

	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/dao/bolt"
	"github.com/cuigh/swirl/dao/mongo"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

// Interface is the interface that wraps all dao methods.
type Interface interface {
	Init() error

	RoleGet(ctx context.Context, id string) (*model.Role, error)
	RoleSearch(ctx context.Context, name string) (roles []*model.Role, err error)
	RoleCreate(ctx context.Context, role *model.Role) error
	RoleUpdate(ctx context.Context, role *model.Role) error
	RoleDelete(ctx context.Context, id string) error

	UserGet(ctx context.Context, id string) (*model.User, error)
	UserGetByName(ctx context.Context, loginName string) (*model.User, error)
	UserSearch(ctx context.Context, args *model.UserSearchArgs) (users []*model.User, count int, err error)
	UserCount(ctx context.Context) (int, error)
	UserCreate(ctx context.Context, user *model.User) error
	UserUpdate(ctx context.Context, user *model.User) error
	UserUpdateStatus(ctx context.Context, user *model.User) error
	UserUpdateProfile(ctx context.Context, user *model.User) error
	UserUpdatePassword(ctx context.Context, user *model.User) error
	UserDelete(ctx context.Context, id string) error

	SessionGet(ctx context.Context, id string) (*model.Session, error)
	SessionCreate(ctx context.Context, session *model.Session) error
	SessionUpdate(ctx context.Context, session *model.Session) error
	SessionUpdateExpiry(ctx context.Context, id string, expiry time.Time) (err error)

	RegistryGet(ctx context.Context, id string) (*model.Registry, error)
	RegistryGetByURL(ctx context.Context, url string) (registry *model.Registry, err error)
	RegistryGetAll(ctx context.Context) (registries []*model.Registry, err error)
	RegistryCreate(ctx context.Context, registry *model.Registry) error
	RegistryUpdate(ctx context.Context, registry *model.Registry) error
	RegistryDelete(ctx context.Context, id string) error

	StackGet(ctx context.Context, name string) (*model.Stack, error)
	StackGetAll(ctx context.Context) (stacks []*model.Stack, err error)
	StackCreate(ctx context.Context, stack *model.Stack) error
	StackUpdate(ctx context.Context, stack *model.Stack) error
	StackDelete(ctx context.Context, name string) error

	EventSearch(ctx context.Context, args *model.EventSearchArgs) (events []*model.Event, count int, err error)
	EventCreate(ctx context.Context, event *model.Event) error

	SettingGet(ctx context.Context, id string) (*model.Setting, error)
	SettingGetAll(ctx context.Context) (settings []*model.Setting, err error)
	SettingUpdate(ctx context.Context, setting *model.Setting) error

	ChartGet(ctx context.Context, id string) (*model.Chart, error)
	ChartGetBatch(ctx context.Context, ids ...string) ([]*model.Chart, error)
	ChartSearch(ctx context.Context, args *model.ChartSearchArgs) (charts []*model.Chart, count int, err error)
	ChartCreate(ctx context.Context, chart *model.Chart) error
	ChartUpdate(ctx context.Context, chart *model.Chart) error
	ChartDelete(ctx context.Context, id string) error

	DashboardGet(ctx context.Context, name, key string) (dashboard *model.Dashboard, err error)
	DashboardUpdate(ctx context.Context, dashboard *model.Dashboard) error
}

func newInterface() (i Interface) {
	var err error

	switch misc.Options.DBType {
	case "", "mongo":
		i, err = mongo.New(misc.Options.DBAddress)
	case "bolt":
		i, err = bolt.New(misc.Options.DBAddress)
	default:
		err = errors.New("Unknown database type: " + misc.Options.DBType)
	}

	if err != nil {
		panic(err)
	}

	return i
}

func init() {
	container.Put(newInterface)
}

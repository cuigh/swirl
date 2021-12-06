package dao

import (
	"context"

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
	Close()

	RoleGet(ctx context.Context, id string) (*model.Role, error)
	RoleList(ctx context.Context, name string) (roles []*model.Role, err error)
	RoleCreate(ctx context.Context, role *model.Role) error
	RoleUpdate(ctx context.Context, role *model.Role) error
	RoleDelete(ctx context.Context, id string) error

	UserCreate(ctx context.Context, user *model.User) error
	UserUpdate(ctx context.Context, user *model.User) error
	UserList(ctx context.Context, args *model.UserSearchArgs) (users []*model.User, count int, err error)
	UserCount(ctx context.Context) (int, error)
	UserGetByID(ctx context.Context, id string) (*model.User, error)
	UserGetByName(ctx context.Context, loginName string) (*model.User, error)
	UserSetStatus(ctx context.Context, id string, status int32) error
	UserDelete(ctx context.Context, id string) error
	UserModifyProfile(ctx context.Context, user *model.User) error
	UserModifyPassword(ctx context.Context, id, pwd, salt string) error

	SessionUpdate(ctx context.Context, session *model.Session) error
	SessionGet(ctx context.Context, token string) (*model.Session, error)

	RegistryCreate(ctx context.Context, registry *model.Registry) error
	RegistryUpdate(ctx context.Context, registry *model.Registry) error
	RegistryGet(ctx context.Context, id string) (*model.Registry, error)
	RegistryGetByURL(ctx context.Context, url string) (registry *model.Registry, err error)
	RegistryList(ctx context.Context) (registries []*model.Registry, err error)
	RegistryDelete(ctx context.Context, id string) error

	StackList(ctx context.Context) (stacks []*model.Stack, err error)
	StackGet(ctx context.Context, name string) (*model.Stack, error)
	StackCreate(ctx context.Context, stack *model.Stack) error
	StackUpdate(ctx context.Context, stack *model.Stack) error
	StackDelete(ctx context.Context, name string) error

	EventCreate(ctx context.Context, event *model.Event) error
	EventList(ctx context.Context, args *model.EventListArgs) (events []*model.Event, count int, err error)

	SettingList(ctx context.Context) (settings []*model.Setting, err error)
	SettingGet(ctx context.Context, id string) (*model.Setting, error)
	SettingUpdate(ctx context.Context, id string, opts []*model.SettingOption) error

	ChartGet(ctx context.Context, id string) (*model.Chart, error)
	ChartBatch(ctx context.Context, ids ...string) ([]*model.Chart, error)
	ChartList(ctx context.Context, title, dashboard string, pageIndex, pageSize int) (charts []*model.Chart, count int, err error)
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

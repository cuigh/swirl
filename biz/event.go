package biz

import (
	"context"
	"time"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/ext/times"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventType string

const (
	EventTypeRegistry  EventType = "Registry"
	EventTypeNode      EventType = "Node"
	EventTypeNetwork   EventType = "Network"
	EventTypeService   EventType = "Service"
	EventTypeStack     EventType = "Stack"
	EventTypeConfig    EventType = "Config"
	EventTypeSecret    EventType = "Secret"
	EventTypeImage     EventType = "Image"
	EventTypeContainer EventType = "Container"
	EventTypeVolume    EventType = "Volume"
	EventTypeUser      EventType = "User"
	EventTypeRole      EventType = "Role"
	EventTypeChart     EventType = "Chart"
	EventTypeSetting   EventType = "Setting"
)

type EventAction string

const (
	EventActionLogin      EventAction = "Login"
	EventActionCreate     EventAction = "Create"
	EventActionDelete     EventAction = "Delete"
	EventActionUpdate     EventAction = "Update"
	EventActionScale      EventAction = "Scale"
	EventActionRollback   EventAction = "Rollback"
	EventActionRestart    EventAction = "Restart"
	EventActionDisconnect EventAction = "Disconnect"
	EventActionDeploy     EventAction = "Deploy"
	EventActionShutdown   EventAction = "Shutdown"
	EventActionPrune      EventAction = "Prune"
)

type EventBiz interface {
	Search(args *dao.EventSearchArgs) (events []*dao.Event, total int, err error)
	Prune(days int32) (err error)
	CreateRegistry(action EventAction, id, name string, user web.User)
	CreateNode(action EventAction, id, name string, user web.User)
	CreateNetwork(action EventAction, id, name string, user web.User)
	CreateService(action EventAction, name string, user web.User)
	CreateConfig(action EventAction, id, name string, user web.User)
	CreateSecret(action EventAction, id, name string, user web.User)
	CreateStack(action EventAction, name string, user web.User)
	CreateImage(action EventAction, node, id string, user web.User)
	CreateContainer(action EventAction, node, id, name string, user web.User)
	CreateVolume(action EventAction, node, name string, user web.User)
	CreateUser(action EventAction, id, name string, user web.User)
	CreateRole(action EventAction, id, name string, user web.User)
	CreateChart(action EventAction, id, title string, user web.User)
	CreateSetting(action EventAction, user web.User)
}

func NewEvent(d dao.Interface) EventBiz {
	return &eventBiz{d: d}
}

type eventBiz struct {
	d dao.Interface
}

func (b *eventBiz) Search(args *dao.EventSearchArgs) (events []*dao.Event, total int, err error) {
	return b.d.EventSearch(context.TODO(), args)
}

func (b *eventBiz) Prune(days int32) (err error) {
	return b.d.EventPrune(context.TODO(), time.Now().Add(-times.Days(days)))
}

func (b *eventBiz) create(et EventType, ea EventAction, args data.Map, user web.User) {
	event := &dao.Event{
		ID:       primitive.NewObjectID(),
		Type:     string(et),
		Action:   string(ea),
		Args:     args,
		UserID:   user.ID(),
		Username: user.Name(),
		Time:     now(),
	}
	err := b.d.EventCreate(context.TODO(), event)
	if err != nil {
		log.Get("event").Errorf("failed to create event `%+v`: %s", event, err)
	}
}

func (b *eventBiz) CreateRegistry(action EventAction, id, name string, user web.User) {
	args := data.Map{"id": id, "name": name}
	b.create(EventTypeRegistry, action, args, user)
}

func (b *eventBiz) CreateService(action EventAction, name string, user web.User) {
	args := data.Map{"name": name}
	b.create(EventTypeService, action, args, user)
}

func (b *eventBiz) CreateNetwork(action EventAction, id, name string, user web.User) {
	args := data.Map{"id": id, "name": name}
	b.create(EventTypeNetwork, action, args, user)
}

func (b *eventBiz) CreateNode(action EventAction, id, name string, user web.User) {
	args := data.Map{"id": id, "name": name}
	b.create(EventTypeNode, action, args, user)
}

func (b *eventBiz) CreateImage(action EventAction, node, id string, user web.User) {
	args := data.Map{"node": node}
	if id != "" {
		args["id"] = id
	}
	b.create(EventTypeImage, action, args, user)
}

func (b *eventBiz) CreateContainer(action EventAction, node, id, name string, user web.User) {
	args := data.Map{"node": node}
	if id != "" {
		args["id"] = id
	}
	if name != "" {
		args["name"] = name
	}
	b.create(EventTypeContainer, action, args, user)
}

func (b *eventBiz) CreateVolume(action EventAction, node, name string, user web.User) {
	args := data.Map{"node": node}
	if name != "" {
		args["name"] = name
	}
	b.create(EventTypeVolume, action, args, user)
}

func (b *eventBiz) CreateStack(action EventAction, name string, user web.User) {
	args := data.Map{"name": name}
	b.create(EventTypeStack, action, args, user)
}

func (b *eventBiz) CreateSecret(action EventAction, id, name string, user web.User) {
	args := data.Map{"id": id, "name": name}
	b.create(EventTypeSecret, action, args, user)
}

func (b *eventBiz) CreateConfig(action EventAction, id, name string, user web.User) {
	args := data.Map{"id": id, "name": name}
	b.create(EventTypeConfig, action, args, user)
}

func (b *eventBiz) CreateRole(action EventAction, id, name string, user web.User) {
	args := data.Map{"id": id, "name": name}
	b.create(EventTypeRole, action, args, user)
}

func (b *eventBiz) CreateSetting(action EventAction, user web.User) {
	b.create(EventTypeSetting, action, nil, user)
}

func (b *eventBiz) CreateUser(action EventAction, id, name string, user web.User) {
	args := data.Map{"id": id, "name": name}
	b.create(EventTypeUser, action, args, user)
}

func (b *eventBiz) CreateChart(action EventAction, id, title string, user web.User) {
	args := data.Map{"id": id, "name": title}
	b.create(EventTypeChart, action, args, user)
}

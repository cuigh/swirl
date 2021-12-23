package biz

import (
	"context"

	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventType string

const (
	EventTypeRegistry        EventType = "Registry"
	EventTypeNode            EventType = "Node"
	EventTypeNetwork         EventType = "Network"
	EventTypeService         EventType = "Service"
	EventTypeServiceTemplate EventType = "Template"
	EventTypeStack           EventType = "Stack"
	EventTypeConfig          EventType = "Config"
	EventTypeSecret          EventType = "Secret"
	EventTypeImage           EventType = "Image"
	EventTypeContainer       EventType = "Container"
	EventTypeVolume          EventType = "Volume"
	EventTypeUser            EventType = "User"
	EventTypeRole            EventType = "Role"
	EventTypeChart           EventType = "Chart"
	EventTypeSetting         EventType = "Setting"
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
	Search(args *model.EventSearchArgs) (events []*model.Event, total int, err error)
	CreateRegistry(action EventAction, id, name string, user web.User)
	CreateNode(action EventAction, id, name string, user web.User)
	CreateNetwork(action EventAction, id, name string, user web.User)
	CreateService(action EventAction, name string, user web.User)
	CreateTemplate(action EventAction, id, name string, user web.User)
	CreateConfig(action EventAction, id, name string, user web.User)
	CreateSecret(action EventAction, id, name string, user web.User)
	CreateStack(action EventAction, name string, user web.User)
	CreateImage(action EventAction, id string, user web.User)
	CreateContainer(action EventAction, id, name string, user web.User)
	CreateVolume(action EventAction, name string, user web.User)
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

func (b *eventBiz) Search(args *model.EventSearchArgs) (events []*model.Event, total int, err error) {
	return b.d.EventSearch(context.TODO(), args)
}

func (b *eventBiz) create(et EventType, ea EventAction, code, name string, user web.User) {
	event := &model.Event{
		ID:       primitive.NewObjectID(),
		Type:     string(et),
		Action:   string(ea),
		Code:     code,
		Name:     name,
		UserID:   user.ID(),
		Username: user.Name(),
		Time:     now(),
	}
	err := b.d.EventCreate(context.TODO(), event)
	if err != nil {
		log.Get("event").Errorf("failed to create event `%+v`: %v", event, err)
	}
}

func (b *eventBiz) CreateRegistry(action EventAction, id, name string, user web.User) {
	b.create(EventTypeRegistry, action, id, name, user)
}

func (b *eventBiz) CreateService(action EventAction, name string, user web.User) {
	b.create(EventTypeService, action, name, name, user)
}

func (b *eventBiz) CreateTemplate(action EventAction, id, name string, user web.User) {
	b.create(EventTypeServiceTemplate, action, id, name, user)
}

func (b *eventBiz) CreateNetwork(action EventAction, id, name string, user web.User) {
	b.create(EventTypeNetwork, action, id, name, user)
}

func (b *eventBiz) CreateNode(action EventAction, id, name string, user web.User) {
	b.create(EventTypeNode, action, id, name, user)
}

func (b *eventBiz) CreateImage(action EventAction, id string, user web.User) {
	b.create(EventTypeImage, action, id, "", user)
}

func (b *eventBiz) CreateContainer(action EventAction, id, name string, user web.User) {
	b.create(EventTypeContainer, action, id, name, user)
}

func (b *eventBiz) CreateVolume(action EventAction, name string, user web.User) {
	b.create(EventTypeVolume, action, name, name, user)
}

func (b *eventBiz) CreateStack(action EventAction, name string, user web.User) {
	b.create(EventTypeStack, action, name, name, user)
}

func (b *eventBiz) CreateSecret(action EventAction, id, name string, user web.User) {
	b.create(EventTypeSecret, action, id, name, user)
}

func (b *eventBiz) CreateConfig(action EventAction, id, name string, user web.User) {
	b.create(EventTypeConfig, action, id, name, user)
}

func (b *eventBiz) CreateRole(action EventAction, id, name string, user web.User) {
	b.create(EventTypeRole, action, id, name, user)
}

func (b *eventBiz) CreateSetting(action EventAction, user web.User) {
	b.create(EventTypeSetting, action, "", "Setting", user)
}

func (b *eventBiz) CreateUser(action EventAction, id, name string, user web.User) {
	b.create(EventTypeUser, action, id, name, user)
}

func (b *eventBiz) CreateChart(action EventAction, id, title string, user web.User) {
	b.create(EventTypeChart, action, id, title, user)
}

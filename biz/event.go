package biz

import (
	"time"

	"github.com/cuigh/auxo/data/guid"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
)

// Event return a event biz instance.
var Event = &eventBiz{}

type eventBiz struct {
}

func (b *eventBiz) Create(event *model.Event) {
	event.ID = guid.New()
	event.Time = time.Now()

	do(func(d dao.Interface) {
		err := d.EventCreate(event)
		if err != nil {
			log.Get("event").Errorf("Create event `%+v` failed: %v", event, err)
		}
	})
	return
}

func (b *eventBiz) CreateRegistry(action model.EventAction, id, name string, user web.User) {
	event := &model.Event{
		Type:     model.EventTypeRegistry,
		Action:   action,
		Code:     id,
		Name:     name,
		UserID:   user.ID(),
		Username: user.Name(),
	}
	b.Create(event)
}

func (b *eventBiz) CreateService(action model.EventAction, name string, user web.User) {
	event := &model.Event{
		Type:     model.EventTypeService,
		Action:   action,
		Code:     name,
		Name:     name,
		UserID:   user.ID(),
		Username: user.Name(),
	}
	b.Create(event)
}

func (b *eventBiz) CreateServiceTemplate(action model.EventAction, id, name string, user web.User) {
	event := &model.Event{
		Type:     model.EventTypeServiceTemplate,
		Action:   action,
		Code:     id,
		Name:     name,
		UserID:   user.ID(),
		Username: user.Name(),
	}
	b.Create(event)
}

func (b *eventBiz) CreateNetwork(action model.EventAction, id, name string, user web.User) {
	event := &model.Event{
		Type:     model.EventTypeNetwork,
		Action:   action,
		Code:     id,
		Name:     name,
		UserID:   user.ID(),
		Username: user.Name(),
	}
	b.Create(event)
}

func (b *eventBiz) CreateVolume(action model.EventAction, name string, user web.User) {
	event := &model.Event{
		Type:     model.EventTypeVolume,
		Action:   action,
		Code:     name,
		Name:     name,
		UserID:   user.ID(),
		Username: user.Name(),
	}
	b.Create(event)
}

func (b *eventBiz) CreateStackTask(action model.EventAction, name string, user web.User) {
	event := &model.Event{
		Type:     model.EventTypeStackTask,
		Action:   action,
		Code:     name,
		Name:     name,
		UserID:   user.ID(),
		Username: user.Name(),
	}
	b.Create(event)
}

func (b *eventBiz) CreateStackArchive(action model.EventAction, id, name string, user web.User) {
	event := &model.Event{
		Type:     model.EventTypeStackArchive,
		Action:   action,
		Code:     id,
		Name:     name,
		UserID:   user.ID(),
		Username: user.Name(),
	}
	b.Create(event)
}

func (b *eventBiz) CreateSecret(action model.EventAction, name string, user web.User) {
	event := &model.Event{
		Type:     model.EventTypeSecret,
		Action:   action,
		Name:     name,
		UserID:   user.ID(),
		Username: user.Name(),
	}
	b.Create(event)
}

func (b *eventBiz) CreateConfig(action model.EventAction, name string, user web.User) {
	event := &model.Event{
		Type:     model.EventTypeConfig,
		Action:   action,
		Name:     name,
		UserID:   user.ID(),
		Username: user.Name(),
	}
	b.Create(event)
}

func (b *eventBiz) CreateRole(action model.EventAction, id, name string, user web.User) {
	event := &model.Event{
		Type:     model.EventTypeRole,
		Action:   action,
		Code:     id,
		Name:     name,
		UserID:   user.ID(),
		Username: user.Name(),
	}
	b.Create(event)
}

func (b *eventBiz) CreateUser(action model.EventAction, loginName, name string, user web.User) {
	event := &model.Event{
		Type:     model.EventTypeUser,
		Action:   action,
		Code:     loginName,
		Name:     name,
		UserID:   user.ID(),
		Username: user.Name(),
	}
	b.Create(event)
}

func (b *eventBiz) CreateSetting(action model.EventAction, user web.User) {
	event := &model.Event{
		Type:     model.EventTypeSetting,
		Action:   action,
		Code:     "",
		Name:     "Setting",
		UserID:   user.ID(),
		Username: user.Name(),
	}
	b.Create(event)
}

func (b *eventBiz) CreateAuthentication(action model.EventAction, userID, loginName, username string) {
	event := &model.Event{
		Type:     model.EventTypeAuthentication,
		Action:   action,
		Code:     loginName,
		Name:     username,
		UserID:   userID,
		Username: username,
	}
	b.Create(event)
}

func (b *eventBiz) List(args *model.EventListArgs) (events []*model.Event, count int, err error) {
	do(func(d dao.Interface) {
		events, count, err = d.EventList(args)
	})
	return
}

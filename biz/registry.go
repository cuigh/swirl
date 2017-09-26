package biz

import (
	"time"

	"github.com/cuigh/auxo/data/guid"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
)

var Registry = &registryBiz{}

type registryBiz struct {
}

func (b *registryBiz) Create(registry *model.Registry, user web.User) (err error) {
	registry.ID = guid.New()
	registry.CreatedAt = time.Now()
	registry.UpdatedAt = registry.CreatedAt

	do(func(d dao.Interface) {
		err = d.RegistryCreate(registry)
		if err == nil {
			Event.CreateRegistry(model.EventActionCreate, registry.ID, registry.Name, user)
		}
	})
	return
}

func (b *registryBiz) Update(registry *model.Registry, user web.User) (err error) {
	do(func(d dao.Interface) {
		err = d.RegistryUpdate(registry)
		if err == nil {
			Event.CreateRegistry(model.EventActionUpdate, registry.ID, registry.Name, user)
		}
	})
	return
}

func (b *registryBiz) List() (registries []*model.Registry, err error) {
	do(func(d dao.Interface) {
		registries, err = d.RegistryList()
	})
	return
}

func (b *registryBiz) Get(id string) (registry *model.Registry, err error) {
	do(func(d dao.Interface) {
		registry, err = d.RegistryGet(id)
	})
	return
}

func (b *registryBiz) Delete(id string, user web.User) (err error) {
	do(func(d dao.Interface) {
		var registry *model.Registry
		registry, err = d.RegistryGet(id)
		if err != nil {
			return
		}

		err = d.RegistryDelete(id)
		if err == nil {
			Event.CreateRegistry(model.EventActionDelete, id, registry.Name, user)
		}
	})
	return
}

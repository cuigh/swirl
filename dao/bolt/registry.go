package bolt

import (
	"context"

	"github.com/cuigh/swirl/model"
)

const Registry = "registry"

func (d *Dao) RegistryCreate(ctx context.Context, registry *model.Registry) (err error) {
	return d.replace(Registry, registry.ID, registry)
}

func (d *Dao) RegistryUpdate(ctx context.Context, registry *model.Registry) (err error) {
	old := &model.Registry{}
	return d.update(Registry, registry.ID, old, func() interface{} {
		registry.CreatedAt = old.CreatedAt
		registry.CreatedBy = old.CreatedBy
		if registry.Password == "" {
			registry.Password = old.Password
		}
		return registry
	})
}

func (d *Dao) RegistryGetAll(ctx context.Context) (registries []*model.Registry, err error) {
	err = d.each(Registry, func(v []byte) error {
		r := &model.Registry{}
		err = decode(v, r)
		if err != nil {
			return err
		}
		registries = append(registries, r)
		return nil
	})
	return
}

func (d *Dao) RegistryGet(ctx context.Context, id string) (registry *model.Registry, err error) {
	registry = &model.Registry{}
	err = d.get(Registry, id, registry)
	if err == ErrNoRecords {
		return nil, nil
	} else if err != nil {
		registry = nil
	}
	return
}

func (d *Dao) RegistryGetByURL(ctx context.Context, url string) (registry *model.Registry, err error) {
	r := &model.Registry{}
	found, err := d.find(Registry, r, func() bool { return r.URL == url })
	if found {
		return r, nil
	}
	return nil, err
}

func (d *Dao) RegistryDelete(ctx context.Context, id string) (err error) {
	return d.delete(Registry, id)
}

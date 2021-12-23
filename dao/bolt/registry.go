package bolt

import (
	"context"

	"github.com/cuigh/swirl/dao"
)

const Registry = "registry"

func (d *Dao) RegistryCreate(ctx context.Context, registry *dao.Registry) (err error) {
	return d.replace(Registry, registry.ID, registry)
}

func (d *Dao) RegistryUpdate(ctx context.Context, registry *dao.Registry) (err error) {
	old := &dao.Registry{}
	return d.update(Registry, registry.ID, old, func() interface{} {
		registry.CreatedAt = old.CreatedAt
		registry.CreatedBy = old.CreatedBy
		if registry.Password == "" {
			registry.Password = old.Password
		}
		return registry
	})
}

func (d *Dao) RegistryGetAll(ctx context.Context) (registries []*dao.Registry, err error) {
	err = d.each(Registry, func(v []byte) error {
		r := &dao.Registry{}
		err = decode(v, r)
		if err != nil {
			return err
		}
		registries = append(registries, r)
		return nil
	})
	return
}

func (d *Dao) RegistryGet(ctx context.Context, id string) (registry *dao.Registry, err error) {
	registry = &dao.Registry{}
	err = d.get(Registry, id, registry)
	if err == ErrNoRecords {
		return nil, nil
	} else if err != nil {
		registry = nil
	}
	return
}

func (d *Dao) RegistryGetByURL(ctx context.Context, url string) (registry *dao.Registry, err error) {
	r := &dao.Registry{}
	found, err := d.find(Registry, r, func() bool { return r.URL == url })
	if found {
		return r, nil
	}
	return nil, err
}

func (d *Dao) RegistryDelete(ctx context.Context, id string) (err error) {
	return d.delete(Registry, id)
}

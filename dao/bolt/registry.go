package bolt

import (
	"context"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/model"
)

func (d *Dao) RegistryCreate(ctx context.Context, registry *model.Registry) (err error) {
	return d.update("registry", registry.ID, registry)
}

func (d *Dao) RegistryUpdate(ctx context.Context, registry *model.Registry) (err error) {
	return d.batch("registry", func(b *bolt.Bucket) error {
		data := b.Get([]byte(registry.ID))
		if data == nil {
			return errors.New("registry not found: " + registry.ID)
		}

		r := &model.Registry{}
		err = json.Unmarshal(data, r)
		if err != nil {
			return err
		}

		r.Name = registry.Name
		r.URL = registry.URL
		r.Username = registry.Username
		if registry.Password != "" {
			r.Password = registry.Password
		}
		r.UpdatedAt = time.Now()
		data, err = json.Marshal(r)
		if err != nil {
			return err
		}

		return b.Put([]byte(registry.ID), data)
	})
}

func (d *Dao) RegistryList(ctx context.Context) (registries []*model.Registry, err error) {
	err = d.each("registry", func(v Value) error {
		r := &model.Registry{}
		err = v.Unmarshal(r)
		if err != nil {
			return err
		}
		registries = append(registries, r)
		return nil
	})
	return
}

func (d *Dao) RegistryGet(ctx context.Context, id string) (registry *model.Registry, err error) {
	var v Value
	v, err = d.get("registry", id)
	if err == nil {
		if v != nil {
			registry = &model.Registry{}
			err = v.Unmarshal(registry)
		}
	}
	return
}

func (d *Dao) RegistryGetByURL(ctx context.Context, url string) (registry *model.Registry, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("registry"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			r := &model.Registry{}
			err = json.Unmarshal(v, r)
			if err != nil {
				return err
			}
			if r.URL == url {
				registry = r
				return nil
			}
		}
		return nil
	})
	return
}

func (d *Dao) RegistryDelete(ctx context.Context, id string) (err error) {
	return d.delete("registry", id)
}

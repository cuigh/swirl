package bolt

import (
	"context"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/model"
)

func (d *Dao) RoleList(ctx context.Context, name string) (roles []*model.Role, err error) {
	err = d.each("role", func(v Value) error {
		role := &model.Role{}
		err = v.Unmarshal(role)
		if err != nil {
			return err
		}

		if matchAny(name, role.Name) {
			roles = append(roles, role)
		}
		return nil
	})
	return
}

func (d *Dao) RoleCreate(ctx context.Context, role *model.Role) (err error) {
	return d.update("role", role.ID, role)
}

func (d *Dao) RoleGet(ctx context.Context, id string) (role *model.Role, err error) {
	var v Value
	v, err = d.get("role", id)
	if err == nil {
		if v != nil {
			role = &model.Role{}
			err = v.Unmarshal(role)
		}
	}
	return
}

func (d *Dao) RoleUpdate(ctx context.Context, role *model.Role) (err error) {
	return d.batch("role", func(b *bolt.Bucket) error {
		data := b.Get([]byte(role.ID))
		if data == nil {
			return errors.New("role not found: " + role.ID)
		}

		r := &model.Role{}
		err = json.Unmarshal(data, r)
		if err != nil {
			return err
		}

		r.Name = role.Name
		r.Description = role.Description
		r.Perms = role.Perms
		r.UpdatedAt = time.Now()
		data, err = json.Marshal(r)
		if err != nil {
			return err
		}

		return b.Put([]byte(role.ID), data)
	})
}

func (d *Dao) RoleDelete(ctx context.Context, id string) (err error) {
	return d.delete("role", id)
}

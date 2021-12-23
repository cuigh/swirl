package bolt

import (
	"context"

	"github.com/cuigh/swirl/dao"
)

const Role = "role"

func (d *Dao) RoleSearch(ctx context.Context, name string) (roles []*dao.Role, err error) {
	err = d.each(Role, func(v []byte) error {
		role := &dao.Role{}
		err = decode(v, role)
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

func (d *Dao) RoleCreate(ctx context.Context, role *dao.Role) (err error) {
	return d.replace(Role, role.ID, role)
}

func (d *Dao) RoleGet(ctx context.Context, id string) (role *dao.Role, err error) {
	role = &dao.Role{}
	err = d.get(Role, id, role)
	if err == ErrNoRecords {
		return nil, nil
	} else if err != nil {
		role = nil
	}
	return
}

func (d *Dao) RoleUpdate(ctx context.Context, role *dao.Role) (err error) {
	old := &dao.Role{}
	return d.update(Role, role.ID, old, func() interface{} {
		role.CreatedAt = old.CreatedAt
		role.CreatedBy = old.CreatedBy
		return role
	})
}

func (d *Dao) RoleDelete(ctx context.Context, id string) (err error) {
	return d.delete(Role, id)
}

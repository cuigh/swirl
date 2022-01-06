package biz

import (
	"context"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
)

type RoleBiz interface {
	Search(ctx context.Context, name string) ([]*dao.Role, error)
	Find(ctx context.Context, id string) (role *dao.Role, err error)
	Create(ctx context.Context, role *dao.Role, user web.User) (err error)
	Delete(ctx context.Context, id, name string, user web.User) (err error)
	Update(ctx context.Context, r *dao.Role, user web.User) (err error)
	GetPerms(ctx context.Context, ids []string) ([]string, error)
}

func NewRole(d dao.Interface, eb EventBiz) RoleBiz {
	return &roleBiz{d: d, eb: eb}
}

type roleBiz struct {
	d  dao.Interface
	eb EventBiz
}

func (b *roleBiz) Search(ctx context.Context, name string) (roles []*dao.Role, err error) {
	return b.d.RoleSearch(ctx, name)
}

func (b *roleBiz) Create(ctx context.Context, role *dao.Role, user web.User) (err error) {
	r := &dao.Role{
		ID:          createId(),
		Name:        role.Name,
		Description: role.Description,
		Perms:       role.Perms,
		CreatedAt:   now(),
		CreatedBy:   newOperator(user),
	}
	r.UpdatedAt = r.CreatedAt
	r.UpdatedBy = r.CreatedBy
	err = b.d.RoleCreate(ctx, r)
	if err == nil {
		b.eb.CreateRole(EventActionCreate, r.ID, role.Name, user)
	}
	return
}

func (b *roleBiz) Delete(ctx context.Context, id, name string, user web.User) (err error) {
	err = b.d.RoleDelete(ctx, id)
	if err == nil {
		go func() {
			_ = b.d.SessionUpdateDirty(ctx, "", id)
			b.eb.CreateRole(EventActionDelete, id, name, user)
		}()
	}
	return
}

func (b *roleBiz) Find(ctx context.Context, id string) (role *dao.Role, err error) {
	return b.d.RoleGet(ctx, id)
}

func (b *roleBiz) Update(ctx context.Context, role *dao.Role, user web.User) (err error) {
	r := &dao.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Perms:       role.Perms,
		UpdatedAt:   now(),
		UpdatedBy:   newOperator(user),
	}
	err = b.d.RoleUpdate(ctx, r)
	if err == nil {
		go func() {
			_ = b.d.SessionUpdateDirty(ctx, "", role.ID)
			b.eb.CreateRole(EventActionUpdate, role.ID, role.Name, user)
		}()
	}
	return
}

func (b *roleBiz) GetPerms(ctx context.Context, ids []string) ([]string, error) {
	m := make(map[string]struct{})

	for _, id := range ids {
		r, err := b.d.RoleGet(ctx, id)
		if err != nil {
			return nil, err
		}

		for _, p := range r.Perms {
			m[p] = struct{}{}
		}
	}

	perms := make([]string, 0, len(m))
	for p := range m {
		perms = append(perms, p)
	}
	return perms, nil
}

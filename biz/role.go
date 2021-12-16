package biz

import (
	"context"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
)

type RoleBiz interface {
	Search(name string) ([]*model.Role, error)
	Find(id string) (role *model.Role, err error)
	Create(role *model.Role, user web.User) (err error)
	Delete(id, name string, user web.User) (err error)
	Update(r *model.Role, user web.User) (err error)
}

func NewRole(d dao.Interface, eb EventBiz) RoleBiz {
	return &roleBiz{d: d, eb: eb}
}

type roleBiz struct {
	d  dao.Interface
	eb EventBiz
}

func (b *roleBiz) Search(name string) (roles []*model.Role, err error) {
	return b.d.RoleSearch(context.TODO(), name)
}

func (b *roleBiz) Create(role *model.Role, user web.User) (err error) {
	r := &model.Role{
		ID:          createId(),
		Name:        role.Name,
		Description: role.Description,
		Perms:       role.Perms,
		CreatedAt:   now(),
		CreatedBy:   newOperator(user),
	}
	r.UpdatedAt = r.CreatedAt
	r.UpdatedBy = r.CreatedBy
	err = b.d.RoleCreate(context.TODO(), r)
	if err == nil {
		b.eb.CreateRole(EventActionCreate, r.ID, role.Name, user)
	}
	return
}

func (b *roleBiz) Delete(id, name string, user web.User) (err error) {
	err = b.d.RoleDelete(context.TODO(), id)
	if err == nil {
		b.eb.CreateRole(EventActionDelete, id, name, user)
	}
	return
}

func (b *roleBiz) Find(id string) (role *model.Role, err error) {
	return b.d.RoleGet(context.TODO(), id)
}

func (b *roleBiz) Update(role *model.Role, user web.User) (err error) {
	r := &model.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Perms:       role.Perms,
		UpdatedAt:   now(),
		UpdatedBy:   newOperator(user),
	}
	err = b.d.RoleUpdate(context.TODO(), r)
	if err == nil {
		b.eb.CreateRole(EventActionUpdate, role.ID, role.Name, user)
	}
	return
}

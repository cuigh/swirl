package biz

import (
	"time"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

// Role return a role biz instance.
var Role = &roleBiz{}

type roleBiz struct {
}

func (b *roleBiz) List() (roles []*model.Role, err error) {
	do(func(d dao.Interface) {
		roles, err = d.RoleList()
	})
	return
}

func (b *roleBiz) Create(role *model.Role, user web.User) (err error) {
	do(func(d dao.Interface) {
		role.ID = misc.NewID()
		role.CreatedAt = time.Now()
		role.UpdatedAt = role.CreatedAt
		err = d.RoleCreate(role)
		if err == nil {
			Event.CreateRole(model.EventActionCreate, role.ID, role.Name, user)
		}
	})
	return
}

func (b *roleBiz) Delete(id string, user web.User) (err error) {
	do(func(d dao.Interface) {
		var role *model.Role
		role, err = d.RoleGet(id)
		if err != nil {
			return
		}

		err = d.RoleDelete(id)
		if err == nil {
			Event.CreateRole(model.EventActionDelete, id, role.Name, user)
		}
	})
	return
}

func (b *roleBiz) Get(id string) (role *model.Role, err error) {
	do(func(d dao.Interface) {
		role, err = d.RoleGet(id)
	})
	return
}

func (b *roleBiz) Update(role *model.Role, user web.User) (err error) {
	do(func(d dao.Interface) {
		role.UpdatedAt = time.Now()
		err = d.RoleUpdate(role)
		if err == nil {
			Event.CreateRole(model.EventActionUpdate, role.ID, role.Name, user)
		}
	})
	return
}

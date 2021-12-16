package biz

import (
	"context"
	"time"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
)

type Role struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty" valid:"required"`
	Description string         `json:"desc,omitempty"`
	Perms       []string       `json:"perms,omitempty"`
	CreatedAt   string         `json:"createdAt,omitempty"`
	UpdatedAt   string         `json:"updatedAt,omitempty"`
	CreatedBy   model.Operator `json:"createdBy" bson:"created_by"`
	UpdatedBy   model.Operator `json:"updatedBy" bson:"updated_by"`
}

func newRole(r *model.Role) *Role {
	return &Role{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Perms:       r.Perms,
		CreatedAt:   formatTime(r.CreatedAt),
		UpdatedAt:   formatTime(r.UpdatedAt),
		CreatedBy:   r.CreatedBy,
		UpdatedBy:   r.UpdatedBy,
	}
}

type RoleBiz interface {
	Search(name string) ([]*Role, error)
	Find(id string) (role *Role, err error)
	Create(role *Role, user web.User) (err error)
	Delete(id, name string, user web.User) (err error)
	Update(r *Role, user web.User) (err error)
}

func NewRole(d dao.Interface, eb EventBiz) RoleBiz {
	return &roleBiz{d: d, eb: eb}
}

type roleBiz struct {
	d  dao.Interface
	eb EventBiz
}

func (b *roleBiz) Search(name string) (roles []*Role, err error) {
	var list []*model.Role
	list, err = b.d.RoleSearch(context.TODO(), name)
	if err == nil {
		for _, r := range list {
			roles = append(roles, newRole(r))
		}
	}
	return
}

func (b *roleBiz) Create(role *Role, user web.User) (err error) {
	r := &model.Role{
		ID:          createId(),
		Name:        role.Name,
		Description: role.Description,
		Perms:       role.Perms,
		CreatedAt:   time.Now(),
	}
	r.UpdatedAt = r.CreatedAt
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

func (b *roleBiz) Find(id string) (role *Role, err error) {
	var r *model.Role
	r, err = b.d.RoleGet(context.TODO(), id)
	if r != nil {
		role = newRole(r)
	}
	return
}

func (b *roleBiz) Update(role *Role, user web.User) (err error) {
	r := &model.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Perms:       role.Perms,
		UpdatedAt:   time.Now(),
	}
	r.UpdatedBy.ID = user.ID()
	r.UpdatedBy.Name = user.Name()
	err = b.d.RoleUpdate(context.TODO(), r)
	if err == nil {
		b.eb.CreateRole(EventActionUpdate, role.ID, role.Name, user)
	}
	return
}

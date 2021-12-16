package biz

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
	"github.com/docker/docker/api/types"
)

type RegistryBiz interface {
	Search() ([]*Registry, error)
	Find(id string) (*Registry, error)
	GetAuth(url string) (auth string, err error)
	Delete(id, name string, user web.User) (err error)
	Create(registry *Registry, user web.User) (err error)
	Update(registry *Registry, user web.User) (err error)
}

func NewRegistry(d dao.Interface, eb EventBiz) RegistryBiz {
	return &registryBiz{d: d, eb: eb}
}

type registryBiz struct {
	d  dao.Interface
	eb EventBiz
}

func (b *registryBiz) Create(registry *Registry, user web.User) (err error) {
	r := registry.Convert()
	r.ID = createId()
	r.CreatedAt = time.Now()
	r.UpdatedAt = r.CreatedAt
	r.CreatedBy.ID = user.ID()
	r.CreatedBy.Name = user.Name()
	r.UpdatedBy = r.CreatedBy

	err = b.d.RegistryCreate(context.TODO(), r)
	if err == nil {
		b.eb.CreateRegistry(EventActionCreate, r.ID, r.Name, user)
	}
	return
}

func (b *registryBiz) Update(registry *Registry, user web.User) (err error) {
	r := registry.Convert()
	r.UpdatedAt = time.Now()
	r.UpdatedBy.ID = user.ID()
	r.UpdatedBy.Name = user.Name()
	err = b.d.RegistryUpdate(context.TODO(), r)
	if err == nil {
		b.eb.CreateRegistry(EventActionUpdate, registry.ID, registry.Name, user)
	}
	return
}

func (b *registryBiz) Search() (registries []*Registry, err error) {
	var list []*model.Registry
	if list, err = b.d.RegistryGetAll(context.TODO()); err == nil {
		for _, r := range list {
			registries = append(registries, newRegistry(r))
		}
	}
	return
}

func (b *registryBiz) Find(id string) (registry *Registry, err error) {
	var r *model.Registry
	if r, err = b.d.RegistryGet(context.TODO(), id); err == nil {
		registry = newRegistry(r)
	}
	return
}

func (b *registryBiz) GetAuth(url string) (auth string, err error) {
	var (
		r   *model.Registry
		buf []byte
	)
	if r, err = b.d.RegistryGetByURL(context.TODO(), url); err == nil && r != nil {
		cfg := &types.AuthConfig{
			ServerAddress: r.URL,
			Username:      r.Username,
			Password:      r.Password,
		}
		if buf, err = json.Marshal(cfg); err == nil {
			auth = base64.URLEncoding.EncodeToString(buf)
		}
	}
	return
}

func (b *registryBiz) Delete(id, name string, user web.User) (err error) {
	err = b.d.RegistryDelete(context.TODO(), id)
	if err == nil {
		b.eb.CreateRegistry(EventActionDelete, id, name, user)
	}
	return
}

type Registry struct {
	ID        string         `json:"id,omitempty"`
	Name      string         `json:"name" valid:"required"`
	URL       string         `json:"url" valid:"required,url"`
	Username  string         `json:"username" valid:"required"`
	Password  string         `json:"password" copier:"-"`
	CreatedAt string         `json:"createdAt,omitempty" copier:"-"`
	UpdatedAt string         `json:"updatedAt,omitempty" copier:"-"`
	CreatedBy model.Operator `json:"createdBy" bson:"created_by"`
	UpdatedBy model.Operator `json:"updatedBy" bson:"updated_by"`
}

func newRegistry(r *model.Registry) *Registry {
	if r == nil {
		return nil
	}

	return &Registry{
		ID:        r.ID,
		Name:      r.Name,
		URL:       r.URL,
		Username:  r.Username,
		CreatedAt: formatTime(r.CreatedAt),
		UpdatedAt: formatTime(r.UpdatedAt),
		CreatedBy: r.CreatedBy,
		UpdatedBy: r.UpdatedBy,
		//Password:  r.Password, // omit password
	}
}

func (r *Registry) Convert() *model.Registry {
	return &model.Registry{
		ID:       r.ID,
		Name:     r.Name,
		URL:      r.URL,
		Username: r.Username,
		Password: r.Password,
	}
}

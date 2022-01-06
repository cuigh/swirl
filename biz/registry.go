package biz

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/docker/docker/api/types"
)

type RegistryBiz interface {
	Search(ctx context.Context) ([]*dao.Registry, error)
	Find(ctx context.Context, id string) (*dao.Registry, error)
	GetAuth(ctx context.Context, url string) (auth string, err error)
	Delete(ctx context.Context, id, name string, user web.User) (err error)
	Create(ctx context.Context, registry *dao.Registry, user web.User) (err error)
	Update(ctx context.Context, registry *dao.Registry, user web.User) (err error)
}

func NewRegistry(d dao.Interface, eb EventBiz) RegistryBiz {
	return &registryBiz{d: d, eb: eb}
}

type registryBiz struct {
	d  dao.Interface
	eb EventBiz
}

func (b *registryBiz) Create(ctx context.Context, r *dao.Registry, user web.User) (err error) {
	r.ID = createId()
	r.CreatedAt = now()
	r.UpdatedAt = r.CreatedAt
	r.CreatedBy = newOperator(user)
	r.UpdatedBy = r.CreatedBy

	err = b.d.RegistryCreate(ctx, r)
	if err == nil {
		b.eb.CreateRegistry(EventActionCreate, r.ID, r.Name, user)
	}
	return
}

func (b *registryBiz) Update(ctx context.Context, r *dao.Registry, user web.User) (err error) {
	r.UpdatedAt = now()
	r.UpdatedBy = newOperator(user)
	err = b.d.RegistryUpdate(ctx, r)
	if err == nil {
		b.eb.CreateRegistry(EventActionUpdate, r.ID, r.Name, user)
	}
	return
}

func (b *registryBiz) Search(ctx context.Context) (registries []*dao.Registry, err error) {
	registries, err = b.d.RegistryGetAll(ctx)
	if err == nil {
		for _, r := range registries {
			r.Password = ""
		}
	}
	return
}

func (b *registryBiz) Find(ctx context.Context, id string) (registry *dao.Registry, err error) {
	registry, err = b.d.RegistryGet(ctx, id)
	if err == nil {
		registry.Password = ""
	}
	return
}

func (b *registryBiz) GetAuth(ctx context.Context, url string) (auth string, err error) {
	var (
		r   *dao.Registry
		buf []byte
	)
	if r, err = b.d.RegistryGetByURL(ctx, url); err == nil && r != nil {
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

func (b *registryBiz) Delete(ctx context.Context, id, name string, user web.User) (err error) {
	err = b.d.RegistryDelete(ctx, id)
	if err == nil {
		b.eb.CreateRegistry(EventActionDelete, id, name, user)
	}
	return
}

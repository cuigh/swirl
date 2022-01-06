package biz

import (
	"context"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/docker"
	"github.com/docker/docker/api/types/swarm"
)

type SecretBiz interface {
	Search(ctx context.Context, name string, pageIndex, pageSize int) (secrets []*Secret, total int, err error)
	Find(ctx context.Context, id string) (secret *Secret, raw string, err error)
	Delete(ctx context.Context, id, name string, user web.User) (err error)
	Create(ctx context.Context, secret *Secret, user web.User) (err error)
	Update(ctx context.Context, secret *Secret, user web.User) (err error)
}

func NewSecret(d *docker.Docker, eb EventBiz) SecretBiz {
	return &secretBiz{d: d, eb: eb}
}

type secretBiz struct {
	d  *docker.Docker
	eb EventBiz
}

func (b *secretBiz) Find(ctx context.Context, id string) (secret *Secret, raw string, err error) {
	var (
		c swarm.Secret
		r []byte
	)
	c, r, err = b.d.SecretInspect(ctx, id)
	if err == nil {
		raw, err = indentJSON(r)
	}
	if err == nil {
		secret = newSecret(&c)
	}
	return
}

func (b *secretBiz) Search(ctx context.Context, name string, pageIndex, pageSize int) ([]*Secret, int, error) {
	list, total, err := b.d.SecretList(ctx, name, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}

	secrets := make([]*Secret, len(list))
	for i, n := range list {
		secrets[i] = newSecret(&n)
	}
	return secrets, total, nil
}

func (b *secretBiz) Delete(ctx context.Context, id, name string, user web.User) (err error) {
	err = b.d.SecretRemove(ctx, id)
	if err == nil {
		b.eb.CreateSecret(EventActionDelete, id, name, user)
	}
	return
}

func (b *secretBiz) Create(ctx context.Context, c *Secret, user web.User) (err error) {
	spec := swarm.SecretSpec{
		Data: []byte(c.Data),
	}
	spec.Name = c.Name
	spec.Labels = toMap(c.Labels)
	if c.Driver.Name != "" {
		spec.Driver = &swarm.Driver{
			Name:    c.Driver.Name,
			Options: toMap(c.Driver.Options),
		}
	}
	if c.Templating.Name != "none" {
		spec.Templating = &swarm.Driver{
			Name:    c.Templating.Name,
			Options: toMap(c.Templating.Options),
		}
	}

	var id string
	id, err = b.d.SecretCreate(ctx, &spec)
	if err != nil {
		b.eb.CreateSecret(EventActionCreate, id, c.Name, user)
	}
	return
}

func (b *secretBiz) Update(ctx context.Context, c *Secret, user web.User) (err error) {
	spec := &swarm.SecretSpec{
		Data: []byte(c.Data),
	}
	spec.Name = c.Name
	spec.Labels = toMap(c.Labels)
	if c.Driver.Name != "" {
		spec.Driver = &swarm.Driver{
			Name:    c.Driver.Name,
			Options: toMap(c.Driver.Options),
		}
	}
	if c.Templating.Name != "" {
		spec.Templating = &swarm.Driver{
			Name:    c.Templating.Name,
			Options: toMap(c.Templating.Options),
		}
	}
	err = b.d.SecretUpdate(ctx, c.ID, c.Version, spec)
	if err == nil {
		b.eb.CreateSecret(EventActionUpdate, c.ID, c.Name, user)
	}
	return
}

type Secret struct {
	ID         string       `json:"id"`
	Name       string       `json:"name,omitempty"`
	Version    uint64       `json:"version"`
	Data       string       `json:"data"`
	Labels     data.Options `json:"labels,omitempty"`
	Driver     Driver       `json:"driver"`
	Templating Driver       `json:"templating"`
	CreatedAt  string       `json:"createdAt"`
	UpdatedAt  string       `json:"updatedAt"`
}

func newSecret(c *swarm.Secret) *Secret {
	secret := &Secret{
		ID:        c.ID,
		Name:      c.Spec.Name,
		Version:   c.Version.Index,
		Data:      string(c.Spec.Data),
		Labels:    mapToOptions(c.Spec.Labels),
		CreatedAt: formatTime(c.CreatedAt),
		UpdatedAt: formatTime(c.UpdatedAt),
	}
	if c.Spec.Driver != nil {
		secret.Driver.Name = c.Spec.Driver.Name
		secret.Driver.Options = mapToOptions(c.Spec.Driver.Options)
	}
	if c.Spec.Templating != nil {
		secret.Templating.Name = c.Spec.Templating.Name
		secret.Templating.Options = mapToOptions(c.Spec.Templating.Options)
	}
	return secret
}

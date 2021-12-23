package biz

import (
	"context"
	"strings"

	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/docker"
	"github.com/cuigh/swirl/docker/compose"
	"github.com/cuigh/swirl/misc"
)

type StackBiz interface {
	Search(name, filter string) (stacks []*dao.Stack, err error)
	Find(name string) (stack *dao.Stack, err error)
	Delete(name string, user web.User) (err error)
	Shutdown(name string, user web.User) (err error)
	Deploy(name string, user web.User) (err error)
	Create(s *dao.Stack, user web.User) (err error)
	Update(s *dao.Stack, user web.User) (err error)
}

func NewStack(d *docker.Docker, s dao.Interface, eb EventBiz) StackBiz {
	return &stackBiz{d: d, s: s, eb: eb}
}

type stackBiz struct {
	d  *docker.Docker
	s  dao.Interface
	eb EventBiz
}

func (b *stackBiz) Search(name, filter string) (stacks []*dao.Stack, err error) {
	var (
		activeStacks   map[string][]string
		internalStacks []*dao.Stack
	)

	// load real stacks
	activeStacks, err = b.d.StackList(context.TODO())
	if err != nil {
		return
	}

	// load stack definitions
	internalStacks, err = b.s.StackGetAll(context.TODO())
	if err != nil {
		return
	}

	// merge stacks and definitions
	for i := range internalStacks {
		stack := internalStacks[i]
		if services, ok := activeStacks[stack.Name]; ok {
			stack.Services = services
			delete(activeStacks, stack.Name)
		}
		if !b.filter(stack, name, filter) {
			stacks = append(stacks, stack)
		}
	}
	for n, services := range activeStacks {
		stack := &dao.Stack{Name: n, Services: services}
		if !b.filter(stack, name, filter) {
			stacks = append(stacks, stack)
		}
	}
	return
}

func (b *stackBiz) Find(name string) (s *dao.Stack, err error) {
	s, err = b.s.StackGet(context.TODO(), name)
	if err != nil {
		return nil, err
	} else if s == nil {
		s = &dao.Stack{Name: name}
	}
	return
}

func (b *stackBiz) filter(stack *dao.Stack, name, filter string) bool {
	if name != "" {
		if !strings.Contains(strings.ToLower(stack.Name), strings.ToLower(name)) {
			return true
		}
	}

	switch filter {
	case "up":
		if len(stack.Services) == 0 {
			return true
		}
	case "internal":
		if !stack.Internal {
			return true
		}
	case "external":
		if stack.Internal {
			return true
		}
	}

	return false
}

func (b *stackBiz) Create(s *dao.Stack, user web.User) (err error) {
	stack := &dao.Stack{
		Name:      s.Name,
		Content:   s.Content,
		CreatedAt: now(),
		CreatedBy: newOperator(user),
	}
	stack.UpdatedAt = stack.CreatedAt
	stack.UpdatedBy = stack.CreatedBy
	err = b.s.StackCreate(context.TODO(), stack)
	if err == nil {
		b.eb.CreateStack(EventActionCreate, stack.Name, user)
	}
	return
}

func (b *stackBiz) Update(s *dao.Stack, user web.User) (err error) {
	stack := &dao.Stack{
		Name:      s.Name,
		Content:   s.Content,
		UpdatedAt: now(),
		UpdatedBy: newOperator(user),
	}
	err = b.s.StackUpdate(context.TODO(), stack)
	if err == nil {
		b.eb.CreateStack(EventActionUpdate, stack.Name, user)
	}
	return
}

func (b *stackBiz) Delete(name string, user web.User) (err error) {
	err = b.d.StackRemove(context.TODO(), name)
	if err == nil {
		err = b.s.StackDelete(context.TODO(), name)
	}
	if err == nil {
		b.eb.CreateStack(EventActionDelete, name, user)
	}
	return
}

func (b *stackBiz) Shutdown(name string, user web.User) (err error) {
	err = b.d.StackRemove(context.TODO(), name)
	if err == nil {
		b.eb.CreateStack(EventActionShutdown, name, user)
	}
	return
}

func (b *stackBiz) Deploy(name string, user web.User) (err error) {
	stack, err := b.s.StackGet(context.TODO(), name)
	if err != nil {
		return err
	} else if stack == nil {
		return errors.Coded(misc.ErrExternalStack, "can not deploy external stack")
	}

	cfg, err := compose.Parse(stack.Name, stack.Content)
	if err != nil {
		return err
	}

	registries, err := b.s.RegistryGetAll(context.TODO())
	if err != nil {
		return err
	}

	// Find auth info from registry
	authes := map[string]string{}
	for _, sc := range cfg.Services {
		if _, ok := authes[sc.Image]; !ok {
			for _, r := range registries {
				if r.Match(sc.Image) {
					authes[sc.Image] = r.GetEncodedAuth()
				}
			}
		}
	}

	err = b.d.StackDeploy(context.TODO(), cfg, authes)
	if err == nil {
		b.eb.CreateStack(EventActionDeploy, name, user)
	}
	return
}

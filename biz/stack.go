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
	"github.com/cuigh/swirl/model"
)

type StackBiz interface {
	Search(name, filter string) (stacks []*Stack, err error)
	Find(name string) (stack *Stack, err error)
	Delete(name string, user web.User) (err error)
	Shutdown(name string, user web.User) (err error)
	Deploy(name string, user web.User) (err error)
	Create(s *Stack, user web.User) (err error)
	Update(s *Stack, user web.User) (err error)
}

func NewStack(d *docker.Docker, s dao.Interface, eb EventBiz) StackBiz {
	return &stackBiz{d: d, s: s, eb: eb}
}

type stackBiz struct {
	d  *docker.Docker
	s  dao.Interface
	eb EventBiz
}

func (b *stackBiz) Search(name, filter string) (stacks []*Stack, err error) {
	var (
		activeStacks   map[string][]string
		internalStacks []*model.Stack
	)

	// load real stacks
	activeStacks, err = b.d.StackList(context.TODO())
	if err != nil {
		return
	}

	// load stack definitions
	internalStacks, err = b.s.StackList(context.TODO())
	if err != nil {
		return
	}

	// merge stacks and definitions
	for _, s := range internalStacks {
		stack := newStack(s)
		if services, ok := activeStacks[stack.Name]; ok {
			stack.Services = services
			delete(activeStacks, stack.Name)
		}
		if !b.filter(stack, name, filter) {
			stacks = append(stacks, stack)
		}
	}
	for n, services := range activeStacks {
		stack := &Stack{Name: n, Services: services}
		if !b.filter(stack, name, filter) {
			stacks = append(stacks, stack)
		}
	}
	return
}

func (b *stackBiz) Find(name string) (stack *Stack, err error) {
	s, err := b.s.StackGet(context.TODO(), name)
	if err != nil {
		return nil, err
	} else if s == nil {
		return &Stack{ID: name, Name: name}, nil
	}
	return newStack(s), nil
}

func (b *stackBiz) filter(stack *Stack, name, filter string) bool {
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

func (b *stackBiz) Create(s *Stack, user web.User) (err error) {
	stack := &model.Stack{
		Name:      s.Name,
		Content:   s.Content,
		CreatedBy: user.ID(),
		UpdatedBy: user.ID(),
	}
	err = b.s.StackCreate(context.TODO(), stack)
	if err == nil {
		b.eb.CreateStack(EventActionCreate, stack.Name, user)
	}
	return
}

func (b *stackBiz) Get(name string) (stack *model.Stack, err error) {
	return b.s.StackGet(context.TODO(), name)
}

func (b *stackBiz) Update(s *Stack, user web.User) (err error) {
	stack := &model.Stack{
		Name:      s.Name,
		Content:   s.Content,
		UpdatedBy: user.ID(),
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

	registries, err := b.s.RegistryList(context.TODO())
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

type Stack struct {
	ID        string   `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Content   string   `json:"content,omitempty"`
	CreatedBy string   `json:"createdBy,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
	UpdatedBy string   `json:"updatedBy,omitempty"`
	UpdatedAt string   `json:"updatedAt,omitempty"`
	Services  []string `json:"services,omitempty"`
	Internal  bool     `json:"internal"`
}

func newStack(s *model.Stack) *Stack {
	return &Stack{
		ID:        s.Name,
		Name:      s.Name,
		Content:   s.Content,
		CreatedBy: s.CreatedBy,
		CreatedAt: formatTime(s.CreatedAt),
		UpdatedBy: s.UpdatedBy,
		UpdatedAt: formatTime(s.UpdatedAt),
		Internal:  true,
	}
}

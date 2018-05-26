package biz

import (
	"strings"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
)

// Stack return a stack biz instance.
var Stack = &stackBiz{}

type stackBiz struct {
}

func (b *stackBiz) List(args *model.StackListArgs) (stacks []*model.Stack, err error) {
	var (
		upStacks, internalStacks []*model.Stack
		upMap                    = make(map[string]*model.Stack)
	)

	// load real stacks
	upStacks, err = docker.StackList()
	if err != nil {
		return
	}
	for _, stack := range upStacks {
		upMap[stack.Name] = stack
	}

	// load stack definitions
	do(func(d dao.Interface) {
		internalStacks, err = d.StackList()
	})
	if err != nil {
		return
	}

	// merge stacks and definitions
	for _, stack := range internalStacks {
		stack.Internal = true
		if s, ok := upMap[stack.Name]; ok {
			stack.Services = s.Services
			delete(upMap, stack.Name)
		}
		if !b.filter(stack, args) {
			stacks = append(stacks, stack)
		}
	}
	for _, stack := range upMap {
		if !b.filter(stack, args) {
			stacks = append(stacks, stack)
		}
	}
	return
}

func (b *stackBiz) filter(stack *model.Stack, args *model.StackListArgs) bool {
	if args.Name != "" {
		if !strings.Contains(strings.ToLower(stack.Name), strings.ToLower(args.Name)) {
			return true
		}
	}

	switch args.Filter {
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

func (b *stackBiz) Create(stack *model.Stack, user web.User) (err error) {
	do(func(d dao.Interface) {
		err = d.StackCreate(stack)
		if err == nil {
			Event.CreateStack(model.EventActionCreate, stack.Name, user)
		}
	})
	return
}

func (b *stackBiz) Get(name string) (stack *model.Stack, err error) {
	do(func(d dao.Interface) {
		stack, err = d.StackGet(name)
	})
	return
}

func (b *stackBiz) Update(stack *model.Stack, user web.User) (err error) {
	do(func(d dao.Interface) {
		err = d.StackUpdate(stack)
		if err == nil {
			Event.CreateStack(model.EventActionUpdate, stack.Name, user)
		}
	})
	return
}

func (b *stackBiz) Delete(name string, user web.User) (err error) {
	do(func(d dao.Interface) {
		err = d.StackDelete(name)
		if err == nil {
			Event.CreateStack(model.EventActionDelete, name, user)
		}
	})
	return
}

// Migrate migrates old archives to stack collection.
func (b *stackBiz) Migrate() {
	do(func(d dao.Interface) {
		d.StackMigrate()
	})
}

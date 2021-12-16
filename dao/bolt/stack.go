package bolt

import (
	"context"

	"github.com/cuigh/swirl/model"
)

const Stack = "stack"

func (d *Dao) StackGetAll(ctx context.Context) (stacks []*model.Stack, err error) {
	err = d.each(Stack, func(v []byte) error {
		stack := &model.Stack{}
		err = decode(v, stack)
		if err == nil {
			stacks = append(stacks, stack)
		}
		return err
	})
	return
}

func (d *Dao) StackCreate(ctx context.Context, stack *model.Stack) (err error) {
	return d.replace(Stack, stack.Name, stack)
}

func (d *Dao) StackGet(ctx context.Context, name string) (stack *model.Stack, err error) {
	stack = &model.Stack{}
	err = d.get(Stack, name, stack)
	if err == ErrNoRecords {
		return nil, nil
	} else if err != nil {
		stack = nil
	}
	return
}

func (d *Dao) StackUpdate(ctx context.Context, stack *model.Stack) (err error) {
	old := &model.Stack{}
	return d.update(Role, stack.Name, old, func() interface{} {
		stack.CreatedAt = old.CreatedAt
		stack.CreatedBy = old.CreatedBy
		return stack
	})
}

func (d *Dao) StackDelete(ctx context.Context, name string) (err error) {
	return d.delete(Stack, name)
}

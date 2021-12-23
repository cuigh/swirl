package mongo

import (
	"context"

	"github.com/cuigh/swirl/dao"
	"go.mongodb.org/mongo-driver/bson"
)

const Stack = "stack"

func (d *Dao) StackGetAll(ctx context.Context) (stacks []*dao.Stack, err error) {
	stacks = []*dao.Stack{}
	err = d.fetch(ctx, Stack, bson.M{}, &stacks)
	return
}

func (d *Dao) StackCreate(ctx context.Context, stack *dao.Stack) (err error) {
	return d.create(ctx, Stack, stack)
}

func (d *Dao) StackGet(ctx context.Context, name string) (stack *dao.Stack, err error) {
	stack = &dao.Stack{}
	found, err := d.find(ctx, Stack, name, stack)
	if !found {
		return nil, err
	}
	return
}

func (d *Dao) StackUpdate(ctx context.Context, stack *dao.Stack) (err error) {
	update := bson.M{
		"$set": bson.M{
			"content":    stack.Content,
			"updated_by": stack.UpdatedBy,
			"updated_at": stack.UpdatedAt,
		},
	}
	return d.update(ctx, Stack, stack.Name, update)
}

func (d *Dao) StackDelete(ctx context.Context, name string) (err error) {
	return d.delete(ctx, Stack, name)
}

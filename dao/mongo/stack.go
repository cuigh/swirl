package mongo

import (
	"context"
	"time"

	"github.com/cuigh/swirl/model"
	"go.mongodb.org/mongo-driver/bson"
)

const Stack = "stack"

func (d *Dao) StackList(ctx context.Context) (stacks []*model.Stack, err error) {
	stacks = []*model.Stack{}
	err = d.fetch(ctx, Stack, bson.M{}, &stacks)
	return
}

func (d *Dao) StackCreate(ctx context.Context, stack *model.Stack) (err error) {
	stack.CreatedAt = time.Now()
	stack.UpdatedAt = stack.CreatedAt
	return d.create(ctx, Stack, stack)
}

func (d *Dao) StackGet(ctx context.Context, name string) (stack *model.Stack, err error) {
	stack = &model.Stack{}
	found, err := d.find(ctx, Stack, name, stack)
	if !found {
		return nil, err
	}
	return
}

func (d *Dao) StackUpdate(ctx context.Context, stack *model.Stack) (err error) {
	update := bson.M{
		"$set": bson.M{
			"content":    stack.Content,
			"updated_by": stack.UpdatedBy,
			"updated_at": time.Now(),
		},
	}
	return d.update(ctx, Stack, stack.Name, update)
}

func (d *Dao) StackDelete(ctx context.Context, name string) (err error) {
	return d.delete(ctx, Stack, name)
}

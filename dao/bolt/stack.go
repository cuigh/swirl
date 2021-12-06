package bolt

import (
	"context"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/model"
)

func (d *Dao) StackList(ctx context.Context) (stacks []*model.Stack, err error) {
	err = d.each("stack", func(v Value) error {
		stack := &model.Stack{}
		err = v.Unmarshal(stack)
		if err == nil {
			stacks = append(stacks, stack)
		}
		return err
	})
	return
}

func (d *Dao) StackCreate(ctx context.Context, stack *model.Stack) (err error) {
	stack.CreatedAt = time.Now()
	stack.UpdatedAt = stack.CreatedAt
	return d.update("stack", stack.Name, stack)
}

func (d *Dao) StackGet(ctx context.Context, name string) (stack *model.Stack, err error) {
	var v Value
	v, err = d.get("stack", name)
	if err == nil {
		if v != nil {
			stack = &model.Stack{}
			err = v.Unmarshal(stack)
		}
	}
	return
}

func (d *Dao) StackUpdate(ctx context.Context, stack *model.Stack) (err error) {
	return d.batch("stack", func(b *bolt.Bucket) error {
		data := b.Get([]byte(stack.Name))
		if data == nil {
			return errors.New("stack not found: " + stack.Name)
		}

		s := &model.Stack{}
		err = json.Unmarshal(data, s)
		if err != nil {
			return err
		}

		s.Content = stack.Content
		s.UpdatedBy = stack.UpdatedBy
		s.UpdatedAt = time.Now()
		data, err = json.Marshal(s)
		if err != nil {
			return err
		}

		return b.Put([]byte(stack.Name), data)
	})
}

func (d *Dao) StackDelete(ctx context.Context, name string) (err error) {
	return d.delete("stack", name)
}

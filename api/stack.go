package api

import (
	"io"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/docker/compose"
	"github.com/cuigh/swirl/misc"
)

// StackHandler encapsulates stack related handlers.
type StackHandler struct {
	Search   web.HandlerFunc `path:"/search" auth:"stack.view" desc:"search stacks"`
	Find     web.HandlerFunc `path:"/find" auth:"stack.view" desc:"find stack by name"`
	Delete   web.HandlerFunc `path:"/delete" method:"post" auth:"stack.delete" desc:"delete stack"`
	Shutdown web.HandlerFunc `path:"/shutdown" method:"post" auth:"stack.shutdown" desc:"shutdown stack"`
	Deploy   web.HandlerFunc `path:"/deploy" method:"post" auth:"stack.deploy" desc:"deploy stack"`
	Save     web.HandlerFunc `path:"/save" method:"post" auth:"stack.edit" desc:"create or update stack"`
	Upload   web.HandlerFunc `path:"/upload" method:"post" auth:"*" desc:"upload compose file"`
}

// NewStack creates an instance of StackHandler
func NewStack(b biz.StackBiz) *StackHandler {
	return &StackHandler{
		Search:   stackSearch(b),
		Find:     stackFind(b),
		Delete:   stackDelete(b),
		Shutdown: stackShutdown(b),
		Deploy:   stackDeploy(b),
		Save:     stackSave(b),
		Upload:   stackUpload,
	}
}

func stackUpload(c web.Context) (err error) {
	file, _, err := c.File("content")
	if err != nil {
		return err
	}
	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return c.Data(b)
}

func stackSearch(b biz.StackBiz) web.HandlerFunc {
	type Args struct {
		Name   string `json:"name" bind:"name"`
		Filter string `json:"filter" bind:"filter"`
	}

	return func(c web.Context) (err error) {
		var (
			args   = &Args{}
			stacks []*dao.Stack
		)

		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			stacks, err = b.Search(ctx, args.Name, args.Filter)
		}

		if err != nil {
			return
		}

		return success(c, stacks)
	}
}

func stackFind(b biz.StackBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		name := c.Query("name")
		stack, err := b.Find(ctx, name)
		if err != nil {
			return err
		}
		return success(c, stack)
	}
}

func stackDelete(b biz.StackBiz) web.HandlerFunc {
	type Args struct {
		Name string `json:"name"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Delete(ctx, args.Name, c.User())
		}
		return ajax(c, err)
	}
}

func stackShutdown(b biz.StackBiz) web.HandlerFunc {
	type Args struct {
		Name string `json:"name"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Shutdown(ctx, args.Name, c.User())
		}
		return ajax(c, err)
	}
}

func stackDeploy(b biz.StackBiz) web.HandlerFunc {
	type Args struct {
		Name string `json:"name"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Deploy(ctx, args.Name, c.User())
		}
		return ajax(c, err)
	}
}

func stackSave(b biz.StackBiz) web.HandlerFunc {
	type Args struct {
		ID string `json:"id"`
		dao.Stack
	}

	return func(c web.Context) error {
		stack := &Args{}
		err := c.Bind(stack, true)
		if err != nil {
			return err
		}

		// Validate format
		_, err = compose.Parse(stack.Name, stack.Content)
		if err != nil {
			return err
		}

		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		if stack.ID == "" {
			err = b.Create(ctx, &stack.Stack, c.User())
		} else {
			err = b.Update(ctx, &stack.Stack, c.User())
		}
		return ajax(c, err)
	}
}

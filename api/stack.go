package api

import (
	"io"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/docker/compose"
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

func stackUpload(ctx web.Context) (err error) {
	file, _, err := ctx.File("content")
	if err != nil {
		return err
	}
	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return ctx.Data(b)
}

func stackSearch(b biz.StackBiz) web.HandlerFunc {
	type Args struct {
		Name   string `json:"name" bind:"name"`
		Filter string `json:"filter" bind:"filter"`
	}

	return func(ctx web.Context) (err error) {
		var (
			args   = &Args{}
			stacks []*dao.Stack
		)

		if err = ctx.Bind(args); err == nil {
			stacks, err = b.Search(args.Name, args.Filter)
		}

		if err != nil {
			return
		}

		return success(ctx, stacks)
	}
}

func stackFind(b biz.StackBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		name := ctx.Query("name")
		stack, err := b.Find(name)
		if err != nil {
			return err
		}
		return success(ctx, stack)
	}
}

func stackDelete(b biz.StackBiz) web.HandlerFunc {
	type Args struct {
		Name string `json:"name"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Delete(args.Name, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func stackShutdown(b biz.StackBiz) web.HandlerFunc {
	type Args struct {
		Name string `json:"name"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Shutdown(args.Name, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func stackDeploy(b biz.StackBiz) web.HandlerFunc {
	type Args struct {
		Name string `json:"name"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Deploy(args.Name, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func stackSave(b biz.StackBiz) web.HandlerFunc {
	type Args struct {
		ID string `json:"id"`
		dao.Stack
	}

	return func(ctx web.Context) error {
		stack := &Args{}
		err := ctx.Bind(stack, true)
		if err != nil {
			return err
		}

		// Validate format
		_, err = compose.Parse(stack.Name, stack.Content)
		if err != nil {
			return err
		}

		if stack.ID == "" {
			err = b.Create(&stack.Stack, ctx.User())
		} else {
			err = b.Update(&stack.Stack, ctx.User())
		}
		return ajax(ctx, err)
	}
}

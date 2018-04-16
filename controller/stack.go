package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/biz/docker/compose"
	"github.com/cuigh/swirl/model"
)

// StackController is a controller of docker stack(compose)
type StackController struct {
	List     web.HandlerFunc `path:"/" name:"stack.list" authorize:"!" desc:"stack list page"`
	New      web.HandlerFunc `path:"/new" name:"stack.new" authorize:"!" desc:"new stack page"`
	Create   web.HandlerFunc `path:"/new" method:"post" name:"stack.create" authorize:"!" desc:"create stack"`
	Detail   web.HandlerFunc `path:"/:name/detail" name:"stack.detail" authorize:"!" desc:"stack detail page"`
	Edit     web.HandlerFunc `path:"/:name/edit" name:"stack.edit" authorize:"!" desc:"stack edit page"`
	Update   web.HandlerFunc `path:"/:name/update" method:"post" name:"stack.update" authorize:"!" desc:"update stack"`
	Deploy   web.HandlerFunc `path:"/:name/deploy" method:"post" name:"stack.deploy" authorize:"!" desc:"deploy stack"`
	Shutdown web.HandlerFunc `path:"/:name/shutdown" method:"post" name:"stack.shutdown" authorize:"!" desc:"shutdown stack"`
	Delete   web.HandlerFunc `path:"/:name/delete" method:"post" name:"stack.delete" authorize:"!" desc:"delete stack"`
}

// Stack creates an instance of StackController
func Stack() (c *StackController) {
	return &StackController{
		List:     stackList,
		New:      stackNew,
		Create:   stackCreate,
		Detail:   stackDetail,
		Edit:     stackEdit,
		Update:   stackUpdate,
		Deploy:   stackDeploy,
		Shutdown: stackShutdown,
		Delete:   stackDelete,
	}
}

func stackList(ctx web.Context) error {
	args := &model.StackListArgs{}
	err := ctx.Bind(args)
	if err != nil {
		return err
	}

	stacks, err := biz.Stack.List(args)
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Stacks", stacks).
		Set("Name", args.Name).
		Set("Filter", args.Filter)
	return ctx.Render("stack/list", m)
}

func stackNew(ctx web.Context) error {
	m := newModel(ctx)
	return ctx.Render("stack/new", m)
}

func stackCreate(ctx web.Context) error {
	stack := &model.Stack{}
	err := ctx.Bind(stack)
	if err != nil {
		return err
	}

	// Validate format
	_, err = compose.Parse(stack.Name, stack.Content)
	if err != nil {
		return err
	}

	stack.CreatedBy = ctx.User().ID()
	stack.UpdatedBy = stack.CreatedBy
	err = biz.Stack.Create(stack, ctx.User())
	return ajaxResult(ctx, err)
}

func stackDetail(ctx web.Context) error {
	name := ctx.P("name")
	stack, err := biz.Stack.Get(name)
	if err != nil {
		return err
	}
	if stack == nil {
		return web.ErrNotFound
	}

	m := newModel(ctx).Set("Stack", stack)
	return ctx.Render("stack/detail", m)
}

func stackEdit(ctx web.Context) error {
	name := ctx.P("name")
	stack, err := biz.Stack.Get(name)
	if err != nil {
		return err
	}
	if stack == nil {
		return web.ErrNotFound
	}

	m := newModel(ctx).Set("Stack", stack)
	return ctx.Render("stack/edit", m)
}

func stackUpdate(ctx web.Context) error {
	stack := &model.Stack{}
	err := ctx.Bind(stack)
	if err == nil {
		// Validate format
		_, err = compose.Parse(stack.Name, stack.Content)
		if err != nil {
			return err
		}

		stack.UpdatedBy = ctx.User().ID()
		err = biz.Stack.Update(stack, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func stackDeploy(ctx web.Context) error {
	name := ctx.P("name")
	stack, err := biz.Stack.Get(name)
	if err != nil {
		return err
	}

	cfg, err := compose.Parse(stack.Name, stack.Content)
	if err != nil {
		return err
	}

	registries, err := biz.Registry.List()
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

	err = docker.StackDeploy(stack.Name, stack.Content, authes)
	if err == nil {
		biz.Event.CreateStack(model.EventActionDeploy, name, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func stackShutdown(ctx web.Context) error {
	name := ctx.P("name")
	err := docker.StackRemove(name)
	if err == nil {
		biz.Event.CreateStack(model.EventActionShutdown, name, ctx.User())
	}
	return ajaxResult(ctx, err)
}

func stackDelete(ctx web.Context) error {
	name := ctx.P("name")
	err := docker.StackRemove(name)
	if err == nil {
		err = biz.Stack.Delete(name, ctx.User())
	}
	return ajaxResult(ctx, err)
}

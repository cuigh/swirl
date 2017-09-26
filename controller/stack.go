package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/biz/docker/compose"
	"github.com/cuigh/swirl/model"
)

type StackController struct {
	TaskList      web.HandlerFunc `path:"/task/" name:"stack.task.list" authorize:"!" desc:"stack task list page"`
	TaskDelete    web.HandlerFunc `path:"/task/delete" method:"post" name:"stack.task.delete" authorize:"!" desc:"delete stack task"`
	ArchiveList   web.HandlerFunc `path:"/archive/" name:"stack.archive.list" authorize:"!" desc:"stack archive list page"`
	ArchiveDetail web.HandlerFunc `path:"/archive/:id/detail" name:"stack.archive.detail" authorize:"!" desc:"stack archive detail page"`
	ArchiveEdit   web.HandlerFunc `path:"/archive/:id/edit" name:"stack.archive.edit" authorize:"!" desc:"stack archive edit page"`
	ArchiveUpdate web.HandlerFunc `path:"/archive/:id/update" method:"post" name:"stack.archive.update" authorize:"!" desc:"update stack archive"`
	ArchiveDelete web.HandlerFunc `path:"/archive/delete" method:"post" name:"stack.archive.delete" authorize:"!" desc:"delete stack archive"`
	ArchiveDeploy web.HandlerFunc `path:"/archive/deploy" method:"post" name:"stack.archive.deploy" authorize:"!" desc:"deploy stack archive"`
	ArchiveNew    web.HandlerFunc `path:"/archive/new" name:"stack.archive.new" authorize:"!" desc:"new stack.archive page"`
	ArchiveCreate web.HandlerFunc `path:"/archive/new" method:"post" name:"stack.archive.create" authorize:"!" desc:"create stack.archive"`
}

func Stack() (c *StackController) {
	c = &StackController{}

	c.TaskList = func(ctx web.Context) error {
		stacks, err := docker.StackList()
		if err != nil {
			return err
		}

		m := newModel(ctx).Add("Stacks", stacks)
		return ctx.Render("stack/task/list", m)
	}

	c.TaskDelete = func(ctx web.Context) error {
		name := ctx.F("name")
		err := docker.StackRemove(name)
		if err == nil {
			biz.Event.CreateStackTask(model.EventActionDelete, name, ctx.User())
		}
		return ajaxResult(ctx, err)
	}

	c.ArchiveList = func(ctx web.Context) error {
		args := &model.ArchiveListArgs{}
		err := ctx.Bind(args)
		if err != nil {
			return err
		}
		args.PageSize = model.PageSize
		if args.PageIndex == 0 {
			args.PageIndex = 1
		}

		archives, totalCount, err := biz.Archive.List(args)
		if err != nil {
			return err
		}

		m := newPagerModel(ctx, totalCount, model.PageSize, args.PageIndex).
			Add("Name", args.Name).
			Add("Archives", archives)
		return ctx.Render("stack/archive/list", m)
	}

	c.ArchiveDetail = func(ctx web.Context) error {
		id := ctx.P("id")
		archive, err := biz.Archive.Get(id)
		if err != nil {
			return err
		}
		if archive == nil {
			return web.ErrNotFound
		}

		m := newModel(ctx).Add("Archive", archive)
		return ctx.Render("stack/archive/detail", m)
	}

	c.ArchiveEdit = func(ctx web.Context) error {
		id := ctx.P("id")
		archive, err := biz.Archive.Get(id)
		if err != nil {
			return err
		}
		if archive == nil {
			return web.ErrNotFound
		}

		m := newModel(ctx).Add("Archive", archive)
		return ctx.Render("stack/archive/edit", m)
	}

	c.ArchiveUpdate = func(ctx web.Context) error {
		archive := &model.Archive{}
		err := ctx.Bind(archive)
		if err == nil {
			// Validate format
			_, err = compose.Parse(archive.Name, archive.Content)
			if err != nil {
				return err
			}

			archive.UpdatedBy = ctx.User().ID()
			err = biz.Archive.Update(archive)
		}
		if err == nil {
			biz.Event.CreateStackArchive(model.EventActionUpdate, archive.ID, archive.Name, ctx.User())
		}
		return ajaxResult(ctx, err)
	}

	c.ArchiveDelete = func(ctx web.Context) error {
		id := ctx.F("id")
		err := biz.Archive.Delete(id, ctx.User())
		return ajaxResult(ctx, err)
	}

	c.ArchiveDeploy = func(ctx web.Context) error {
		id := ctx.F("id")
		archive, err := biz.Archive.Get(id)
		if err != nil {
			return err
		}

		cfg, err := compose.Parse(archive.Name, archive.Content)
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

		err = docker.StackDeploy(archive.Name, archive.Content, authes)
		return ajaxResult(ctx, err)
	}

	c.ArchiveNew = func(ctx web.Context) error {
		m := newModel(ctx)
		return ctx.Render("stack/archive/new", m)
	}

	c.ArchiveCreate = func(ctx web.Context) error {
		archive := &model.Archive{}
		err := ctx.Bind(archive)
		if err == nil {
			// Validate format
			_, err = compose.Parse(archive.Name, archive.Content)
			if err != nil {
				return err
			}

			archive.CreatedBy = ctx.User().ID()
			archive.UpdatedBy = archive.CreatedBy
			err = biz.Archive.Create(archive)
		}
		return ajaxResult(ctx, err)
	}

	return
}

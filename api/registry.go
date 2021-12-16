package api

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

// RegistryHandler encapsulates registry related handlers.
type RegistryHandler struct {
	Search web.HandlerFunc `path:"/search" auth:"registry.view" desc:"search registries"`
	Find   web.HandlerFunc `path:"/find" auth:"registry.view" desc:"find registry by id"`
	Delete web.HandlerFunc `path:"/delete" method:"post" auth:"registry.delete" desc:"delete registry"`
	Save   web.HandlerFunc `path:"/save" method:"post" auth:"registry.edit" desc:"create or update registry"`
}

// NewRegistry creates an instance of RegistryHandler
func NewRegistry(b biz.RegistryBiz) *RegistryHandler {
	return &RegistryHandler{
		Search: registrySearch(b),
		Find:   registryFind(b),
		Delete: registryDelete(b),
		Save:   registrySave(b),
	}
}

func registrySearch(b biz.RegistryBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		registries, err := b.Search()
		if err != nil {
			return err
		}
		return success(ctx, registries)
	}
}

func registryFind(b biz.RegistryBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		id := ctx.Query("id")
		node, err := b.Find(id)
		if err != nil {
			return err
		}
		return success(ctx, node)
	}
}

func registryDelete(b biz.RegistryBiz) web.HandlerFunc {
	type Args struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Delete(args.ID, args.Name, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func registrySave(b biz.RegistryBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		r := &model.Registry{}
		err := ctx.Bind(r, true)
		if err == nil {
			if r.ID == "" {
				err = b.Create(r, ctx.User())
			} else {
				err = b.Update(r, ctx.User())
			}
		}
		return ajax(ctx, err)
	}
}

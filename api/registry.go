package api

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
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
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		registries, err := b.Search(ctx)
		if err != nil {
			return err
		}
		return success(c, registries)
	}
}

func registryFind(b biz.RegistryBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		id := c.Query("id")
		node, err := b.Find(ctx, id)
		if err != nil {
			return err
		}
		return success(c, node)
	}
}

func registryDelete(b biz.RegistryBiz) web.HandlerFunc {
	type Args struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Delete(ctx, args.ID, args.Name, c.User())
		}
		return ajax(c, err)
	}
}

func registrySave(b biz.RegistryBiz) web.HandlerFunc {
	return func(c web.Context) error {
		r := &dao.Registry{}
		err := c.Bind(r, true)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			if r.ID == "" {
				err = b.Create(ctx, r, c.User())
			} else {
				err = b.Update(ctx, r, c.User())
			}
		}
		return ajax(c, err)
	}
}

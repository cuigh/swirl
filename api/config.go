package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
)

// ConfigHandler encapsulates config related handlers.
type ConfigHandler struct {
	Search web.HandlerFunc `path:"/search" auth:"config.view" desc:"search configs"`
	Find   web.HandlerFunc `path:"/find" auth:"config.view" desc:"find config by name"`
	Delete web.HandlerFunc `path:"/delete" method:"post" auth:"config.delete" desc:"delete config"`
	Save   web.HandlerFunc `path:"/save" method:"post" auth:"config.edit" desc:"create or update config"`
}

// NewConfig creates an instance of ConfigHandler
func NewConfig(b biz.ConfigBiz) *ConfigHandler {
	return &ConfigHandler{
		Search: configSearch(b),
		Find:   configFind(b),
		Delete: configDelete(b),
		Save:   configSave(b),
	}
}

func configSearch(b biz.ConfigBiz) web.HandlerFunc {
	type Args struct {
		Name      string `json:"name" bind:"name"`
		PageIndex int    `json:"pageIndex" bind:"pageIndex"`
		PageSize  int    `json:"pageSize" bind:"pageSize"`
	}

	return func(ctx web.Context) (err error) {
		var (
			args    = &Args{}
			configs []*biz.Config
			total   int
		)

		if err = ctx.Bind(args); err == nil {
			configs, total, err = b.Search(args.Name, args.PageIndex, args.PageSize)
		}

		if err != nil {
			return
		}

		return success(ctx, data.Map{
			"items": configs,
			"total": total,
		})
	}
}

func configFind(b biz.ConfigBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		id := ctx.Query("id")
		config, raw, err := b.Find(id)
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"config": config, "raw": raw})
	}
}

func configDelete(b biz.ConfigBiz) web.HandlerFunc {
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

func configSave(b biz.ConfigBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		c := &biz.Config{}
		err := ctx.Bind(c, true)
		if err == nil {
			if c.ID == "" {
				err = b.Create(c, ctx.User())
			} else {
				err = b.Update(c, ctx.User())
			}
		}
		return ajax(ctx, err)
	}
}

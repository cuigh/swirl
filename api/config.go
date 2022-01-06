package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/misc"
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

	return func(c web.Context) (err error) {
		var (
			args    = &Args{}
			configs []*biz.Config
			total   int
		)

		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			configs, total, err = b.Search(ctx, args.Name, args.PageIndex, args.PageSize)
		}

		if err != nil {
			return
		}

		return success(c, data.Map{
			"items": configs,
			"total": total,
		})
	}
}

func configFind(b biz.ConfigBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		id := c.Query("id")
		config, raw, err := b.Find(ctx, id)
		if err != nil {
			return err
		}
		return success(c, data.Map{"config": config, "raw": raw})
	}
}

func configDelete(b biz.ConfigBiz) web.HandlerFunc {
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

func configSave(b biz.ConfigBiz) web.HandlerFunc {
	return func(c web.Context) error {
		config := &biz.Config{}
		err := c.Bind(config, true)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			if config.ID == "" {
				err = b.Create(ctx, config, c.User())
			} else {
				err = b.Update(ctx, config, c.User())
			}
		}
		return ajax(c, err)
	}
}

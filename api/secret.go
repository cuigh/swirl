package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
)

// SecretHandler encapsulates secret related handlers.
type SecretHandler struct {
	Search web.HandlerFunc `path:"/search" auth:"secret.view" desc:"search secrets"`
	Find   web.HandlerFunc `path:"/find" auth:"secret.view" desc:"find secret by name"`
	Delete web.HandlerFunc `path:"/delete" method:"post" auth:"secret.delete" desc:"delete secret"`
	Save   web.HandlerFunc `path:"/save" method:"post" auth:"secret.edit" desc:"create or update secret"`
}

// NewSecret creates an instance of SecretHandler
func NewSecret(b biz.SecretBiz) *SecretHandler {
	return &SecretHandler{
		Search: secretSearch(b),
		Find:   secretFind(b),
		Delete: secretDelete(b),
		Save:   secretSave(b),
	}
}

func secretSearch(b biz.SecretBiz) web.HandlerFunc {
	type Args struct {
		Name      string `json:"name" bind:"name"`
		PageIndex int    `json:"pageIndex" bind:"pageIndex"`
		PageSize  int    `json:"pageSize" bind:"pageSize"`
	}

	return func(ctx web.Context) (err error) {
		var (
			args    = &Args{}
			secrets []*biz.Secret
			total   int
		)

		if err = ctx.Bind(args); err == nil {
			secrets, total, err = b.Search(args.Name, args.PageIndex, args.PageSize)
		}

		if err != nil {
			return
		}

		return success(ctx, data.Map{
			"items": secrets,
			"total": total,
		})
	}
}

func secretFind(b biz.SecretBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		id := ctx.Query("id")
		secret, raw, err := b.Find(id)
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"secret": secret, "raw": raw})
	}
}

func secretDelete(b biz.SecretBiz) web.HandlerFunc {
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

func secretSave(b biz.SecretBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		c := &biz.Secret{}
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

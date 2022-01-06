package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/misc"
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

	return func(c web.Context) (err error) {
		var (
			args    = &Args{}
			secrets []*biz.Secret
			total   int
		)

		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			secrets, total, err = b.Search(ctx, args.Name, args.PageIndex, args.PageSize)
		}

		if err != nil {
			return
		}

		return success(c, data.Map{
			"items": secrets,
			"total": total,
		})
	}
}

func secretFind(b biz.SecretBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		id := c.Query("id")
		secret, raw, err := b.Find(ctx, id)
		if err != nil {
			return err
		}
		return success(c, data.Map{"secret": secret, "raw": raw})
	}
}

func secretDelete(b biz.SecretBiz) web.HandlerFunc {
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

func secretSave(b biz.SecretBiz) web.HandlerFunc {
	return func(c web.Context) error {
		secret := &biz.Secret{}
		err := c.Bind(secret, true)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			if secret.ID == "" {
				err = b.Create(ctx, secret, c.User())
			} else {
				err = b.Update(ctx, secret, c.User())
			}
		}
		return ajax(c, err)
	}
}

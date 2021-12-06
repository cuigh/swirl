package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
)

// ImageHandler encapsulates image related handlers.
type ImageHandler struct {
	Search     web.HandlerFunc `path:"/search" auth:"image.view" desc:"search images"`
	Find       web.HandlerFunc `path:"/find" auth:"image.view" desc:"find image by id"`
	Delete     web.HandlerFunc `path:"/delete" method:"post" auth:"image.delete" desc:"delete image"`
}

// NewImage creates an instance of ImageHandler
func NewImage(b biz.ImageBiz) *ImageHandler {
	return &ImageHandler{
		Search:     imageSearch(b),
		Find:       imageFind(b),
		Delete:     imageDelete(b),
	}
}

func imageSearch(b biz.ImageBiz) web.HandlerFunc {
	type Args struct {
		Name      string `json:"name" bind:"name"`
		PageIndex int    `json:"pageIndex" bind:"pageIndex"`
		PageSize  int    `json:"pageSize" bind:"pageSize"`
	}

	return func(ctx web.Context) (err error) {
		var (
			args   = &Args{}
			images []*biz.Image
			total  int
		)

		if err = ctx.Bind(args); err == nil {
			images, total, err = b.Search(args.Name, args.PageIndex, args.PageSize)
		}

		if err != nil {
			return
		}

		return success(ctx, data.Map{
			"items": images,
			"total": total,
		})
	}
}

func imageFind(b biz.ImageBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		id := ctx.Query("id")
		image, raw, err := b.Find(id)
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"image": image, "raw": raw})
	}
}

func imageDelete(b biz.ImageBiz) web.HandlerFunc {
	type Args struct {
		ID   string `json:"id"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Delete(args.ID, ctx.User())
		}
		return ajax(ctx, err)
	}
}

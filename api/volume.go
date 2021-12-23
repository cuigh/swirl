package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
)

// VolumeHandler encapsulates volume related handlers.
type VolumeHandler struct {
	Search web.HandlerFunc `path:"/search" auth:"volume.view" desc:"search volumes"`
	Find   web.HandlerFunc `path:"/find" auth:"volume.view" desc:"find volume by name"`
	Delete web.HandlerFunc `path:"/delete" method:"post" auth:"volume.delete" desc:"delete volume"`
	Save   web.HandlerFunc `path:"/save" method:"post" auth:"volume.edit" desc:"create or update volume"`
	Prune  web.HandlerFunc `path:"/prune" method:"post" auth:"volume.delete" desc:"delete unused volumes"`
}

// NewVolume creates an instance of VolumeHandler
func NewVolume(b biz.VolumeBiz) *VolumeHandler {
	return &VolumeHandler{
		Search: volumeSearch(b),
		Find:   volumeFind(b),
		Delete: volumeDelete(b),
		Save:   volumeSave(b),
		Prune:  volumePrune(b),
	}
}

func volumeSearch(b biz.VolumeBiz) web.HandlerFunc {
	type Args struct {
		Node      string `json:"node" bind:"node"`
		Name      string `json:"name" bind:"name"`
		PageIndex int    `json:"pageIndex" bind:"pageIndex"`
		PageSize  int    `json:"pageSize" bind:"pageSize"`
	}

	return func(ctx web.Context) (err error) {
		var (
			args    = &Args{}
			volumes []*biz.Volume
			total   int
		)

		if err = ctx.Bind(args); err == nil {
			volumes, total, err = b.Search(args.Node, args.Name, args.PageIndex, args.PageSize)
		}

		if err != nil {
			return
		}

		return success(ctx, data.Map{
			"items": volumes,
			"total": total,
		})
	}
}

func volumeFind(b biz.VolumeBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		node := ctx.Query("node")
		name := ctx.Query("name")
		volume, raw, err := b.Find(node, name)
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"volume": volume, "raw": raw})
	}
}

func volumeDelete(b biz.VolumeBiz) web.HandlerFunc {
	type Args struct {
		Node string `json:"node"`
		Name string `json:"name"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Delete(args.Node, args.Name, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func volumeSave(b biz.VolumeBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		v := &biz.Volume{}
		err := ctx.Bind(v, true)
		if err == nil {
			err = b.Create(v, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func volumePrune(b biz.VolumeBiz) web.HandlerFunc {
	type Args struct {
		Node string `json:"node"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err != nil {
			return err
		}

		count, size, err := b.Prune(args.Node, ctx.User())
		if err != nil {
			return err
		}

		return success(ctx, data.Map{
			"count": count,
			"size":  size,
		})
	}
}

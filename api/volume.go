package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/misc"
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

	return func(c web.Context) (err error) {
		var (
			args    = &Args{}
			volumes []*biz.Volume
			total   int
		)

		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			volumes, total, err = b.Search(ctx, args.Node, args.Name, args.PageIndex, args.PageSize)
		}

		if err != nil {
			return
		}

		return success(c, data.Map{
			"items": volumes,
			"total": total,
		})
	}
}

func volumeFind(b biz.VolumeBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		node := c.Query("node")
		name := c.Query("name")
		volume, raw, err := b.Find(ctx, node, name)
		if err != nil {
			return err
		}
		return success(c, data.Map{"volume": volume, "raw": raw})
	}
}

func volumeDelete(b biz.VolumeBiz) web.HandlerFunc {
	type Args struct {
		Node string `json:"node"`
		Name string `json:"name"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Delete(ctx, args.Node, args.Name, c.User())
		}
		return ajax(c, err)
	}
}

func volumeSave(b biz.VolumeBiz) web.HandlerFunc {
	return func(c web.Context) error {
		v := &biz.Volume{}
		err := c.Bind(v, true)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Create(ctx, v, c.User())
		}
		return ajax(c, err)
	}
}

func volumePrune(b biz.VolumeBiz) web.HandlerFunc {
	type Args struct {
		Node string `json:"node"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err != nil {
			return err
		}

		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		count, size, err := b.Prune(ctx, args.Node, c.User())
		if err != nil {
			return err
		}

		return success(c, data.Map{
			"count": count,
			"size":  size,
		})
	}
}

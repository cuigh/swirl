package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
)

// NetworkHandler encapsulates network related handlers.
type NetworkHandler struct {
	Search     web.HandlerFunc `path:"/search" auth:"network.view" desc:"search networks"`
	Find       web.HandlerFunc `path:"/find" auth:"network.view" desc:"find network by name"`
	Delete     web.HandlerFunc `path:"/delete" method:"post" auth:"network.delete" desc:"delete network"`
	Save       web.HandlerFunc `path:"/save" method:"post" auth:"network.edit" desc:"create or update network"`
	Disconnect web.HandlerFunc `path:"/disconnect" method:"post" auth:"network.disconnect" desc:"disconnect container from network"`
}

// NewNetwork creates an instance of NetworkHandler
func NewNetwork(nb biz.NetworkBiz) *NetworkHandler {
	return &NetworkHandler{
		Search:     networkSearch(nb),
		Find:       networkFind(nb),
		Delete:     networkDelete(nb),
		Save:       networkSave(nb),
		Disconnect: networkDisconnect(nb),
	}
}

func networkSearch(nb biz.NetworkBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		networks, err := nb.Search()
		if err != nil {
			return err
		}
		return success(ctx, networks)
	}
}

func networkFind(nb biz.NetworkBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		name := ctx.Query("name")
		network, raw, err := nb.Find(name)
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"network": network, "raw": raw})
	}
}

func networkDelete(nb biz.NetworkBiz) web.HandlerFunc {
	type Args struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = nb.Delete(args.ID, args.Name, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func networkSave(nb biz.NetworkBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		n := &biz.Network{}
		err := ctx.Bind(n, true)
		if err == nil {
			err = nb.Create(n, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func networkDisconnect(nb biz.NetworkBiz) web.HandlerFunc {
	type Args struct {
		NetworkID   string `json:"networkId"`
		NetworkName string `json:"networkName"`
		Container   string `json:"container"`
	}
	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args, true)
		if err == nil {
			err = nb.Disconnect(args.NetworkID, args.NetworkName, args.Container, ctx.User())
		}
		return ajax(ctx, err)
	}
}

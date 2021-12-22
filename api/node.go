package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz"
)

// NodeHandler encapsulates node related handlers.
type NodeHandler struct {
	List   web.HandlerFunc `path:"/list" auth:"node.view" desc:"list nodes"`
	Search web.HandlerFunc `path:"/search" auth:"node.view" desc:"search nodes"`
	Find   web.HandlerFunc `path:"/find" auth:"node.view" desc:"find node by name"`
	Delete web.HandlerFunc `path:"/delete" method:"post" auth:"node.delete" desc:"delete node"`
	Save   web.HandlerFunc `path:"/save" method:"post" auth:"node.edit" desc:"create or update node"`
}

// NewNode creates an instance of NodeHandler
func NewNode(nb biz.NodeBiz) *NodeHandler {
	return &NodeHandler{
		List:   nodeList(nb),
		Search: nodeSearch(nb),
		Find:   nodeFind(nb),
		Delete: nodeDelete(nb),
		Save:   nodeSave(nb),
	}
}

func nodeList(nb biz.NodeBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		agent := cast.ToBool(ctx.Query("agent"))
		nodes, err := nb.List(agent)
		if err != nil {
			return err
		}
		return success(ctx, nodes)
	}
}

func nodeSearch(nb biz.NodeBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		nodes, err := nb.Search()
		if err != nil {
			return err
		}

		return success(ctx, nodes)
	}
}

func nodeFind(nb biz.NodeBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		id := ctx.Query("id")
		node, raw, err := nb.Find(id)
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"node": node, "raw": raw})
	}
}

func nodeDelete(nb biz.NodeBiz) web.HandlerFunc {
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

func nodeSave(nb biz.NodeBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		n := &biz.Node{}
		err := ctx.Bind(n, true)
		if err == nil {
			err = nb.Update(n, ctx.User())
		}
		return ajax(ctx, err)
	}
}

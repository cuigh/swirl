package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/misc"
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
	return func(c web.Context) error {
		agent := cast.ToBool(c.Query("agent"))
		nodes, err := nb.List(agent)
		if err != nil {
			return err
		}
		return success(c, nodes)
	}
}

func nodeSearch(nb biz.NodeBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		nodes, err := nb.Search(ctx)
		if err != nil {
			return err
		}

		return success(c, nodes)
	}
}

func nodeFind(nb biz.NodeBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		id := c.Query("id")
		node, raw, err := nb.Find(ctx, id)
		if err != nil {
			return err
		}
		return success(c, data.Map{"node": node, "raw": raw})
	}
}

func nodeDelete(nb biz.NodeBiz) web.HandlerFunc {
	type Args struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = nb.Delete(ctx, args.ID, args.Name, c.User())
		}
		return ajax(c, err)
	}
}

func nodeSave(nb biz.NodeBiz) web.HandlerFunc {
	return func(c web.Context) error {
		n := &biz.Node{}
		err := c.Bind(n, true)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = nb.Update(ctx, n, c.User())
		}
		return ajax(c, err)
	}
}

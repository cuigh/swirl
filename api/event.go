package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/dao"
)

// EventHandler encapsulates event related handlers.
type EventHandler struct {
	Search web.HandlerFunc `path:"/search" auth:"event.view" desc:"search events"`
	Prune  web.HandlerFunc `path:"/prune" method:"post" auth:"event.delete" desc:"delete events"`
}

// NewEvent creates an instance of EventHandler
func NewEvent(b biz.EventBiz) *EventHandler {
	return &EventHandler{
		Search: eventSearch(b),
		Prune:  eventPrune(b),
	}
}

func eventSearch(b biz.EventBiz) web.HandlerFunc {
	return func(ctx web.Context) (err error) {
		var (
			args   = &dao.EventSearchArgs{}
			events []*dao.Event
			total  int
		)

		if err = ctx.Bind(args); err == nil {
			events, total, err = b.Search(args)
		}

		if err != nil {
			return
		}

		return success(ctx, data.Map{
			"items": events,
			"total": total,
		})
	}
}

func eventPrune(b biz.EventBiz) web.HandlerFunc {
	type Args struct {
		Days int32 `json:"days"`
	}

	return func(ctx web.Context) (err error) {
		var args = &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Prune(args.Days)
		}
		return ajax(ctx, err)
	}
}

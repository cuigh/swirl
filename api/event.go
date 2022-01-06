package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
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
	return func(c web.Context) (err error) {
		var (
			args   = &dao.EventSearchArgs{}
			events []*dao.Event
			total  int
		)

		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			events, total, err = b.Search(ctx, args)
		}

		if err != nil {
			return
		}

		return success(c, data.Map{
			"items": events,
			"total": total,
		})
	}
}

func eventPrune(b biz.EventBiz) web.HandlerFunc {
	type Args struct {
		Days int32 `json:"days"`
	}

	return func(c web.Context) (err error) {
		var args = &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Prune(ctx, args.Days)
		}
		return ajax(c, err)
	}
}

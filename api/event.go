package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
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
			args   = &model.EventSearchArgs{}
			events []*biz.Event
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
		Date string `json:"date"`
	}

	return func(ctx web.Context) (err error) {
		var args = &Args{}
		if err = ctx.Bind(args); err == nil {
			// TODO
			//err = b.Prune(args.Date)
			err = errors.NotImplemented
		}
		return ajax(ctx, err)
	}
}

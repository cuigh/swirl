package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

// EventController is a controller of user events
type EventController struct {
	List web.HandlerFunc `path:"/" name:"event.list" authorize:"!" desc:"event list page"`
}

// Event creates an instance of EventController
func Event() (c *EventController) {
	return &EventController{
		List: eventList,
	}
}

func eventList(ctx web.Context) error {
	args := &model.EventListArgs{}
	err := ctx.Bind(args)
	if err != nil {
		return err
	}
	args.PageSize = model.PageSize
	if args.PageIndex == 0 {
		args.PageIndex = 1
	}

	events, totalCount, err := biz.Event.List(args)
	if err != nil {
		return err
	}

	m := newPagerModel(ctx, totalCount, model.PageSize, args.PageIndex).
		Set("Events", events).Set("Args", args)
	return ctx.Render("system/event/list", m)
}

package controller

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

type EventController struct {
	List web.HandlerFunc `path:"/" name:"event.list" authorize:"!" desc:"event list page"`
}

func Event() (c *EventController) {
	c = &EventController{}

	c.List = func(ctx web.Context) error {
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
			Add("Events", events).Add("Args", args)
		return ctx.Render("system/event/list", m)
	}

	return
}

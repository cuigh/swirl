package bolt

import (
	"context"
	"sort"

	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

func (d *Dao) EventList(ctx context.Context, args *model.EventListArgs) (events []*model.Event, count int, err error) {
	err = d.each("event", func(v Value) error {
		event := &model.Event{}
		err = v.Unmarshal(event)
		if err != nil {
			return err
		}

		match := true
		if args.Name != "" {
			match = matchAny(args.Name, event.Name)
		}
		if match && args.Type != "" {
			match = string(event.Type) == args.Type
		}

		if match {
			events = append(events, event)
		}
		return nil
	})
	if err == nil {
		count = len(events)
		sort.Slice(events, func(i, j int) bool {
			return events[i].Time.After(events[j].Time)
		})
		start, end := misc.Page(count, args.PageIndex, args.PageSize)
		events = events[start:end]
	}
	return
}

func (d *Dao) EventCreate(ctx context.Context, event *model.Event) (err error) {
	return d.update("event", event.ID.Hex(), event)
}

package bolt

import (
	"context"
	"sort"
	"time"

	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

const Event = "event"

func (d *Dao) EventSearch(ctx context.Context, args *model.EventSearchArgs) (events []*model.Event, count int, err error) {
	err = d.each(Event, func(v []byte) error {
		event := &model.Event{}
		err = decode(v, event)
		if err != nil {
			return err
		}

		match := true
		if args.Name != "" {
			match = matchAny(args.Name, event.Name)
		}
		if match && args.Type != "" {
			match = event.Type == args.Type
		}

		if match {
			events = append(events, event)
		}
		return nil
	})
	if err == nil {
		count = len(events)
		sort.Slice(events, func(i, j int) bool {
			return time.Time(events[i].Time).After(time.Time(events[j].Time))
		})
		start, end := misc.Page(count, args.PageIndex, args.PageSize)
		events = events[start:end]
	}
	return
}

func (d *Dao) EventCreate(ctx context.Context, event *model.Event) (err error) {
	return d.replace(Event, event.ID.Hex(), event)
}

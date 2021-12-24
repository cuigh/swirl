package bolt

import (
	"context"
	"sort"
	"time"

	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
)

const Event = "event"

func (d *Dao) EventSearch(ctx context.Context, args *dao.EventSearchArgs) (events []*dao.Event, count int, err error) {
	err = d.each(Event, func(v []byte) error {
		event := &dao.Event{}
		err = decode(v, event)
		if err != nil {
			return err
		}

		match := true
		if args.Name != "" {
			match = event.Args != nil && matchAny(args.Name, cast.ToString(event.Args["name"]))
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

func (d *Dao) EventCreate(ctx context.Context, event *dao.Event) (err error) {
	return d.replace(Event, event.ID.Hex(), event)
}

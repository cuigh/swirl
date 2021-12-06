package mongo

import (
	"context"

	"github.com/cuigh/swirl/model"
	"go.mongodb.org/mongo-driver/bson"
)

const Event = "event"

func (d *Dao) EventList(ctx context.Context, args *model.EventListArgs) (events []*model.Event, count int, err error) {
	filter := bson.M{}
	if args.Type != "" {
		filter["type"] = args.Type
	}
	if args.Name != "" {
		filter["name"] = args.Name
	}
	opts := searchOptions{filter: filter, sorter: bson.M{"_id": -1}, pageIndex: args.PageIndex, pageSize: args.PageSize}
	events = []*model.Event{}
	count, err = d.search(ctx, Event, opts, &events)
	return
}

func (d *Dao) EventCreate(ctx context.Context, event *model.Event) (err error) {
	return d.create(ctx, Event, event)
}

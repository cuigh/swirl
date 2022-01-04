package mongo

import (
	"context"
	"time"

	"github.com/cuigh/swirl/dao"
	"go.mongodb.org/mongo-driver/bson"
)

const Event = "event"

func (d *Dao) EventSearch(ctx context.Context, args *dao.EventSearchArgs) (events []*dao.Event, count int, err error) {
	filter := bson.M{}
	if args.Type != "" {
		filter["type"] = args.Type
	}
	if args.Name != "" {
		filter["args.name"] = args.Name
	}
	opts := searchOptions{filter: filter, sorter: bson.M{"_id": -1}, pageIndex: args.PageIndex, pageSize: args.PageSize}
	events = []*dao.Event{}
	count, err = d.search(ctx, Event, opts, &events)
	return
}

func (d *Dao) EventCreate(ctx context.Context, event *dao.Event) (err error) {
	return d.create(ctx, Event, event)
}

func (d *Dao) EventPrune(ctx context.Context, end time.Time) (err error) {
	filter := bson.M{"time": bson.M{"$lt": end}}
	_, err = d.db.Collection(Event).DeleteMany(ctx, filter)
	return
}

package mongo

import (
	"github.com/cuigh/swirl/model"
	"gopkg.in/mgo.v2/bson"
)

func (d *Dao) EventList(args *model.EventListArgs) (events []*model.Event, count int, err error) {
	d.do(func(db *database) {
		query := bson.M{}
		if args.Type != "" {
			query["type"] = args.Type
		}
		if args.Name != "" {
			query["name"] = args.Name
		}

		q := db.C("event").Find(query)
		count, err = q.Count()
		if err != nil {
			return
		}

		events = []*model.Event{}
		err = q.Sort("-time").Skip(args.PageSize * (args.PageIndex - 1)).Limit(args.PageSize).All(&events)
	})
	return
}

func (d *Dao) EventCreate(event *model.Event) (err error) {
	d.do(func(db *database) {
		err = db.C("event").Insert(event)
	})
	return
}

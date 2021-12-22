package mongo

import (
	"context"
	"time"

	"github.com/cuigh/swirl/model"
	"go.mongodb.org/mongo-driver/bson"
)

const Session = "session"

func (d *Dao) SessionGet(ctx context.Context, id string) (session *model.Session, err error) {
	session = &model.Session{}
	found, err := d.find(ctx, Session, id, session)
	if !found {
		return nil, err
	}
	return
}

func (d *Dao) SessionCreate(ctx context.Context, session *model.Session) (err error) {
	return d.create(ctx, Session, session)
}

func (d *Dao) SessionUpdate(ctx context.Context, session *model.Session) (err error) {
	return d.update(ctx, Session, session.ID, session)
}

func (d *Dao) SessionUpdateExpiry(ctx context.Context, id string, expiry time.Time) (err error) {
	update := bson.M{
		"$set": bson.M{
			"expiry":     expiry,
			"updated_by": time.Now(),
		},
	}
	return d.update(ctx, Session, id, update)
}

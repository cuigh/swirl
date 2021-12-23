package mongo

import (
	"context"
	"time"

	"github.com/cuigh/swirl/dao"
	"go.mongodb.org/mongo-driver/bson"
)

const Session = "session"

func (d *Dao) SessionGet(ctx context.Context, id string) (session *dao.Session, err error) {
	session = &dao.Session{}
	found, err := d.find(ctx, Session, id, session)
	if !found {
		return nil, err
	}
	return
}

func (d *Dao) SessionCreate(ctx context.Context, session *dao.Session) (err error) {
	return d.create(ctx, Session, session)
}

func (d *Dao) SessionUpdate(ctx context.Context, session *dao.Session) (err error) {
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

func (d *Dao) SessionUpdateDirty(ctx context.Context, userID string, roleID string) (err error) {
	filter := bson.M{}
	if userID != "" {
		filter["userId"] = userID
	} else if roleID != "" {
		filter["roles"] = roleID
	} else {
		return nil
	}

	update := bson.M{
		"dirty":      true,
		"updated_by": time.Now(),
	}
	_, err = d.db.Collection(Session).UpdateMany(ctx, filter, update)
	return
}

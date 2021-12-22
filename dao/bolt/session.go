package bolt

import (
	"context"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cuigh/swirl/model"
)

const Session = "session"

func (d *Dao) SessionGet(ctx context.Context, id string) (session *model.Session, err error) {
	s := &model.Session{}
	err = d.get(Session, id, s)
	if err == ErrNoRecords {
		return nil, nil
	} else if err != nil {
		s = nil
	}
	return
}

func (d *Dao) SessionCreate(ctx context.Context, session *model.Session) (err error) {
	return d.replace(Session, session.ID, session)
}

func (d *Dao) SessionUpdate(ctx context.Context, session *model.Session) (err error) {
	return d.replace(Session, session.UserID, session)
}

func (d *Dao) SessionUpdateExpiry(ctx context.Context, id string, expiry time.Time) (err error) {
	old := &model.Session{}
	return d.update(Session, id, old, func() interface{} {
		old.Expiry = expiry
		old.UpdatedAt = time.Now()
		return old
	})
}

// SessionPrune cleans up expired logs.
func (d *Dao) SessionPrune() {
	err := d.db.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(Session))
		return b.ForEach(func(k, v []byte) error {
			session := &model.Session{}
			if err = decode(v, session); err == nil && session.Expiry.Add(time.Hour).Before(time.Now()) {
				err = b.Delete(k)
			}
			return err
		})
	})
	if err != nil {
		d.logger.Error("failed to clean up expired sessions: ", err)
	}
}

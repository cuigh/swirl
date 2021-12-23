package bolt

import (
	"context"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cuigh/swirl/dao"
)

const Session = "session"

func (d *Dao) SessionGet(ctx context.Context, id string) (session *dao.Session, err error) {
	s := &dao.Session{}
	err = d.get(Session, id, s)
	if err == ErrNoRecords {
		return nil, nil
	} else if err != nil {
		s = nil
	}
	return
}

func (d *Dao) SessionCreate(ctx context.Context, session *dao.Session) (err error) {
	return d.replace(Session, session.ID, session)
}

func (d *Dao) SessionUpdate(ctx context.Context, session *dao.Session) (err error) {
	return d.replace(Session, session.UserID, session)
}

func (d *Dao) SessionUpdateExpiry(ctx context.Context, id string, expiry time.Time) (err error) {
	old := &dao.Session{}
	return d.update(Session, id, old, func() interface{} {
		old.Expiry = expiry
		old.UpdatedAt = time.Now()
		return old
	})
}

func (d *Dao) SessionUpdateDirty(ctx context.Context, userID string, roleID string) (err error) {
	contains := func(arr []string, str string) bool {
		for _, s := range arr {
			if s == str {
				return true
			}
		}
		return false
	}

	var (
		buf []byte
		now = time.Now()
	)
	return d.db.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(Session))
		return b.ForEach(func(k, v []byte) error {
			session := &dao.Session{}
			if err = decode(v, session); err != nil {
				return err
			}

			if (userID != "" && session.UserID == userID) || (roleID != "" && contains(session.Roles, roleID)) {
				session.Dirty = true
				session.UpdatedAt = now
				if buf, err = encode(session); err == nil {
					err = b.Put(k, buf)
				}
			}
			return err
		})
	})
}

// SessionPrune cleans up expired logs.
func (d *Dao) SessionPrune() {
	err := d.db.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(Session))
		return b.ForEach(func(k, v []byte) error {
			session := &dao.Session{}
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

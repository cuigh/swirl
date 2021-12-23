package bolt

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/util/run"
	"github.com/cuigh/swirl/dao"
	"go.mongodb.org/mongo-driver/bson"
)

var ErrNoRecords = errors.New("no records")

func encode(v interface{}) ([]byte, error) {
	return bson.Marshal(v)
}

func decode(d []byte, v interface{}) error {
	return bson.Unmarshal(d, v)
}

// Dao implements dao.Interface interface.
type Dao struct {
	db     *bolt.DB
	logger log.Logger
}

// New creates a Dao instance.
func New(addr string) (dao.Interface, error) {
	if addr == "" {
		addr = "/data/swirl"
	}

	db, err := bolt.Open(filepath.Join(addr, "swirl.db"), 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open bolt database")
	}

	d := &Dao{
		logger: log.Get("bolt"),
		db:     db,
	}
	run.Schedule(time.Hour, d.SessionPrune, func(err interface{}) {
		d.logger.Error("failed to clean up expired sessions: ", err)
	})
	return d, nil
}

func (d *Dao) Init() error {
	buckets := []string{"chart", "dashboard", "event", "registry", "role", "setting", "stack", "user" /*"perm","session","template"*/}
	return d.db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *Dao) replace(bucket, key string, value interface{}) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		buf, err := encode(value)
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte(bucket))
		return b.Put([]byte(key), buf)
	})
}

func (d *Dao) update(bucket, key string, oldValue interface{}, newValue func() interface{}) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get([]byte(key))
		if data == nil {
			return ErrNoRecords
		}

		if err := decode(data, oldValue); err != nil {
			return err
		}

		buf, err := encode(newValue())
		if err != nil {
			return err
		}
		return b.Put([]byte(key), buf)
	})
}

func (d *Dao) delete(bucket, key string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete([]byte(key))
	})
}

func (d *Dao) get(bucket, key string, value interface{}) error {
	return d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			if data := b.Get([]byte(key)); data != nil {
				return decode(data, value)
			}
		}
		return ErrNoRecords
	})
}

func (d Dao) find(bucket string, value interface{}, matcher func() bool) (found bool, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucket)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err = decode(v, value); err != nil {
				return err
			}
			if matcher() {
				found = true
				return nil
			}
		}
		return nil
	})
	return
}

func (d *Dao) count(bucket string) (count int, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(bucket)); b != nil {
			count = b.Stats().KeyN
		}
		return nil
	})
	return
}

func (d *Dao) each(bucket string, fn func(v []byte) error) (err error) {
	return d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.ForEach(func(k, v []byte) error {
			return fn(v)
		})
	})
}

func (d *Dao) slice(bucket string, fn func(v []byte) error, keys ...string) (err error) {
	return d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		for _, key := range keys {
			if data := b.Get([]byte(key)); data != nil {
				if err = fn(data); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func matchAny(s string, list ...string) bool {
	s = strings.ToLower(s)
	for _, v := range list {
		if strings.Contains(strings.ToLower(v), s) {
			return true
		}
	}
	return false
}

func init() {
	dao.Register("bolt", New)
}

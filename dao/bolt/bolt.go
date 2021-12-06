package bolt

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/log"
)

type Value []byte

func (v Value) Unmarshal(i interface{}) error {
	return json.Unmarshal([]byte(v), i)
}

// Dao implements dao.Interface interface.
type Dao struct {
	db     *bolt.DB
	logger log.Logger
}

// New creates a Dao instance.
func New(addr string) (*Dao, error) {
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

func (d *Dao) Close() {
	d.db.Close()
}

func (d *Dao) update(bucket, key string, value interface{}) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		buf, err := json.Marshal(value)
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte(bucket))
		return b.Put([]byte(key), buf)
	})
}

func (d *Dao) delete(bucket, key string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete([]byte(key))
	})
}

func (d *Dao) get(bucket, key string) (val Value, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(bucket)); b != nil {
			v := b.Get([]byte(key))
			if v != nil {
				val = Value(v)
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

func (d *Dao) each(bucket string, fn func(v Value) error) (err error) {
	return d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.ForEach(func(k, v []byte) error {
			return fn(Value(v))
		})
	})
}

func (d *Dao) slice(bucket string, fn func(v Value) error, keys ...string) (err error) {
	return d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		for _, key := range keys {
			if data := b.Get([]byte(key)); data != nil {
				if err = fn(Value(data)); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (d *Dao) batch(bucket string, fn func(b *bolt.Bucket) error) (err error) {
	return d.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return fn(b)
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

// itob returns an 8-byte big endian representation of v.
//func itob(i uint64) []byte {
//	b := make([]byte, 8)
//	binary.BigEndian.PutUint64(b, i)
//	return b
//}

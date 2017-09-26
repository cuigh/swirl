package mongo

import (
	"errors"

	"github.com/cuigh/auxo/log"
	"gopkg.in/mgo.v2"
)

var (
	indexes = map[string][]mgo.Index{
		"user": {
			mgo.Index{Key: []string{"login_name"}, Unique: true},
			mgo.Index{Key: []string{"name"}},
			mgo.Index{Key: []string{"email"}, Unique: true},
			mgo.Index{Key: []string{"admin"}},
			mgo.Index{Key: []string{"status"}},
		},
		"role": {
			mgo.Index{Key: []string{"name"}, Unique: true},
		},
		"session": {
			mgo.Index{Key: []string{"token"}, Unique: true},
		},
		"archive": {
			mgo.Index{Key: []string{"name"}, Unique: true},
		},
		"event": {
			mgo.Index{Key: []string{"type"}},
			mgo.Index{Key: []string{"name"}},
			mgo.Index{Key: []string{"-time"}},
		},
	}
)

type database struct {
	db *mgo.Database
}

func (d *database) Close() {
	d.db.Session.Close()
}

func (d *database) C(name string) *mgo.Collection {
	return d.db.C(name)
}

func (d *database) Run(cmd, result interface{}) error {
	return d.db.Run(cmd, result)
}

type Dao struct {
	logger  *log.Logger
	session *mgo.Session
}

func New(addr string) (*Dao, error) {
	if addr == "" {
		return nil, errors.New("database address must be configured for mongo storage")
	}

	s, err := mgo.Dial(addr)
	if err != nil {
		return nil, err
	}

	d := &Dao{
		session: s,
		logger:  log.Get("mongo"),
	}
	d.createIndexes()
	return d, nil
}

func (d *Dao) createIndexes() {
	db := d.db()
	defer db.Close()

	for name, ins := range indexes {
		c := db.C(name)
		for _, in := range ins {
			err := c.EnsureIndex(in)
			if err != nil {
				d.logger.Warnf("Ensure index %s-%v failed: %v", name, in.Key, err)
			}
		}
	}
}

func (d *Dao) db() *database {
	return &database{
		db: d.session.Copy().DB(""),
	}
}

func (d *Dao) do(fn func(db *database)) {
	db := d.db()
	defer db.Close()

	fn(db)
}

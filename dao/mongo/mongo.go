package mongo

import (
	"context"
	"strings"
	"time"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var indexes = map[string][]mongo.IndexModel{
	"chart": {
		mongo.IndexModel{Keys: bson.D{{"title", 1}}},
	},
	"user": {
		mongo.IndexModel{
			Keys:    bson.D{{"login_name", 1}},
			Options: options.Index().SetUnique(true),
		},
		mongo.IndexModel{Keys: bson.D{{"name", 1}}},
	},
	"role": {
		mongo.IndexModel{
			Keys:    bson.D{{"name", 1}},
			Options: options.Index().SetUnique(true),
		},
	},
	"event": {
		mongo.IndexModel{Keys: bson.D{{"type", 1}}},
		mongo.IndexModel{Keys: bson.D{{"name", 1}}},
	},
	//"session": {
	//	mongo.IndexModel{
	//		Keys:    bson.D{{"token", 1}},
	//		Options: options.Index().SetUnique(true),
	//	},
	//},
}

type Dao struct {
	db     *mongo.Database
	logger log.Logger
}

func New(addr string) (*Dao, error) {
	db, err := open(addr)
	if err != nil {
		return nil, err
	}

	return &Dao{
		db:     db,
		logger: log.Get("mongo"),
	}, nil
}

func open(addr string) (*mongo.Database, error) {
	db := "swirl"

	// compatible with old mgo driver
	if !strings.HasPrefix(addr, "mongodb://") && !strings.HasPrefix(addr, "mongodb+srv://") {
		addr = "mongodb://" + addr
	}

	cs, err := connstring.Parse(addr)
	if err != nil {
		return nil, err
	} else if cs.Database != "" {
		db = cs.Database
	}

	opts := &options.ClientOptions{}
	opts.ApplyURI(addr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts.SetAppName(app.Name))
	if err != nil {
		return nil, err
	}
	return client.Database(db), nil
}

func (d *Dao) Init() (err error) {
	for name, models := range indexes {
		c := d.db.Collection(name)
		_, err = c.Indexes().CreateMany(context.TODO(), models)
		if err != nil {
			return
		}
	}
	return
}

func (d *Dao) find(ctx context.Context, coll string, id interface{}, v interface{}) (found bool, err error) {
	err = d.db.Collection(coll).FindOne(ctx, bson.M{"_id": id}).Decode(v)
	if err == nil {
		found = true
	} else if err == mongo.ErrNoDocuments {
		err = nil
	}
	return
}

func (d *Dao) fetch(ctx context.Context, coll string, filter bson.M, records interface{}) (err error) {
	var cur *mongo.Cursor
	cur, err = d.db.Collection(coll).Find(ctx, filter)
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	return cur.All(ctx, records)
}

func (d *Dao) search(ctx context.Context, coll string, opts searchOptions, records interface{}) (int, error) {
	var (
		c   = d.db.Collection(coll)
		cur *mongo.Cursor
	)

	// fetch total count
	count, err := c.CountDocuments(ctx, opts.filter)
	if err != nil {
		return 0, err
	}

	// fetch records
	findOpts := options.Find().SetSkip(int64(opts.pageSize * (opts.pageIndex - 1))).SetLimit(int64(opts.pageSize))
	if opts.sorter != nil {
		findOpts.SetSort(opts.sorter)
	}
	cur, err = c.Find(context.TODO(), opts.filter, findOpts)
	if err != nil {
		return 0, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, records)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (d *Dao) create(ctx context.Context, coll string, doc interface{}) (err error) {
	_, err = d.db.Collection(coll).InsertOne(ctx, doc)
	return
}

func (d *Dao) update(ctx context.Context, coll string, id interface{}, update interface{}) (err error) {
	_, err = d.db.Collection(coll).UpdateByID(ctx, id, update)
	return
}

func (d *Dao) upsert(ctx context.Context, coll string, id interface{}, update interface{}) (err error) {
	_, err = d.db.Collection(coll).UpdateByID(ctx, id, update, options.Update().SetUpsert(true))
	return
}

func (d *Dao) delete(ctx context.Context, coll string, id interface{}) (err error) {
	_, err = d.db.Collection(coll).DeleteOne(ctx, bson.M{"_id": id})
	return
}

type searchOptions struct {
	filter    bson.M
	sorter    bson.M
	pageIndex int
	pageSize  int
}

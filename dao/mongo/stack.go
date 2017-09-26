package mongo

import (
	"time"

	"github.com/cuigh/auxo/data/guid"
	"github.com/cuigh/swirl/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (d *Dao) ArchiveList(args *model.ArchiveListArgs) (archives []*model.Archive, count int, err error) {
	d.do(func(db *database) {
		var query bson.M
		if args.Name != "" {
			query = bson.M{"name": args.Name}
		}
		q := db.C("archive").Find(query)

		count, err = q.Count()
		if err != nil {
			return
		}

		archives = []*model.Archive{}
		err = q.Skip(args.PageSize * (args.PageIndex - 1)).Limit(args.PageSize).All(&archives)
	})
	return
}

func (d *Dao) ArchiveCreate(archive *model.Archive) (err error) {
	archive.ID = guid.New()
	archive.CreatedAt = time.Now()
	archive.UpdatedAt = archive.CreatedAt

	d.do(func(db *database) {
		err = db.C("archive").Insert(archive)
	})
	return
}

func (d *Dao) ArchiveGet(id string) (archive *model.Archive, err error) {
	d.do(func(db *database) {
		archive = &model.Archive{}
		err = db.C("archive").FindId(id).One(archive)
		if err == mgo.ErrNotFound {
			archive, err = nil, nil
		} else if err != nil {
			archive = nil
		}
	})
	return
}

func (d *Dao) ArchiveUpdate(archive *model.Archive) (err error) {
	d.do(func(db *database) {
		update := bson.M{
			"$set": bson.M{
				"name":       archive.Name,
				"content":    archive.Content,
				"updated_by": archive.UpdatedBy,
				"updated_at": time.Now(),
			},
		}
		err = db.C("archive").UpdateId(archive.ID, update)
	})
	return
}

func (d *Dao) ArchiveDelete(id string) (err error) {
	d.do(func(db *database) {
		err = db.C("archive").RemoveId(id)
	})
	return
}

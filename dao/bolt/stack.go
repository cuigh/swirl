package bolt

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/model"
)

//func (d *Dao) ArchiveList(args *model.ArchiveListArgs) (archives []*model.Archive, count int, err error) {
//	d.do(func(db *database) {
//		var query bson.M
//		if args.Name != "" {
//			query = bson.M{"name": args.Name}
//		}
//		q := db.C("archive").Find(query)
//
//		count, err = q.Count()
//		if err != nil {
//			return
//		}
//
//		archives = []*model.Archive{}
//		err = q.Skip(args.PageSize * (args.PageIndex - 1)).Limit(args.PageSize).All(&archives)
//	})
//	return
//}
//
//func (d *Dao) ArchiveCreate(archive *model.Archive) (err error) {
//	archive.ID = misc.NewID()
//	archive.CreatedAt = time.Now()
//	archive.UpdatedAt = archive.CreatedAt
//
//	d.do(func(db *database) {
//		err = db.C("archive").Insert(archive)
//	})
//	return
//}
//
//func (d *Dao) ArchiveGet(id string) (archive *model.Archive, err error) {
//	d.do(func(db *database) {
//		archive = &model.Archive{}
//		err = db.C("archive").FindId(id).One(archive)
//		if err == mgo.ErrNotFound {
//			archive, err = nil, nil
//		} else if err != nil {
//			archive = nil
//		}
//	})
//	return
//}
//
//func (d *Dao) ArchiveUpdate(archive *model.Archive) (err error) {
//	d.do(func(db *database) {
//		update := bson.M{
//			"$set": bson.M{
//				"name":       archive.Name,
//				"content":    archive.Content,
//				"updated_by": archive.UpdatedBy,
//				"updated_at": time.Now(),
//			},
//		}
//		err = db.C("archive").UpdateId(archive.ID, update)
//	})
//	return
//}
//
//func (d *Dao) ArchiveDelete(id string) (err error) {
//	d.do(func(db *database) {
//		err = db.C("archive").RemoveId(id)
//	})
//	return
//}

//===============================

func (d *Dao) StackList() (stacks []*model.Stack, err error) {
	err = d.each("stack", func(v Value) error {
		stack := &model.Stack{}
		err = v.Unmarshal(stack)
		if err == nil {
			stacks = append(stacks, stack)
		}
		return err
	})
	return
}

func (d *Dao) StackCreate(stack *model.Stack) (err error) {
	stack.CreatedAt = time.Now()
	stack.UpdatedAt = stack.CreatedAt
	return d.update("stack", stack.Name, stack)
}

func (d *Dao) StackGet(name string) (stack *model.Stack, err error) {
	var v Value
	v, err = d.get("stack", name)
	if err == nil {
		if v != nil {
			stack = &model.Stack{}
			err = v.Unmarshal(stack)
		}
	}
	return
}

func (d *Dao) StackUpdate(stack *model.Stack) (err error) {
	return d.batch("stack", func(b *bolt.Bucket) error {
		data := b.Get([]byte(stack.Name))
		if data == nil {
			return errors.New("stack not found: " + stack.Name)
		}

		s := &model.Stack{}
		err = json.Unmarshal(data, s)
		if err != nil {
			return err
		}

		s.Content = stack.Content
		s.UpdatedBy = stack.UpdatedBy
		s.UpdatedAt = time.Now()
		data, err = json.Marshal(s)
		if err != nil {
			return err
		}

		return b.Put([]byte(stack.Name), data)
	})
}

func (d *Dao) StackDelete(name string) (err error) {
	return d.delete("stack", name)
}

// StackMigrate migrates stacks from old archive collection.
func (d *Dao) StackMigrate() {
	// bolt storage engine was implemented at version 0.7.8, so migration is not required.
}

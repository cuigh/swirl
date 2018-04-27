package mongo

import (
	"time"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/swirl/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (d *Dao) StackList() (stacks []*model.Stack, err error) {
	d.do(func(db *database) {
		stacks = []*model.Stack{}
		err = db.C("stack").Find(nil).All(&stacks)
	})
	return
}

func (d *Dao) StackCreate(stack *model.Stack) (err error) {
	stack.CreatedAt = time.Now()
	stack.UpdatedAt = stack.CreatedAt

	d.do(func(db *database) {
		err = db.C("stack").Insert(stack)
	})
	return
}

func (d *Dao) StackGet(name string) (stack *model.Stack, err error) {
	d.do(func(db *database) {
		stack = &model.Stack{}
		err = db.C("stack").FindId(name).One(stack)
		if err == mgo.ErrNotFound {
			stack, err = nil, nil
		} else if err != nil {
			stack = nil
		}
	})
	return
}

func (d *Dao) StackUpdate(stack *model.Stack) (err error) {
	d.do(func(db *database) {
		update := bson.M{
			"$set": bson.M{
				"content":    stack.Content,
				"updated_by": stack.UpdatedBy,
				"updated_at": time.Now(),
			},
		}
		err = db.C("stack").UpdateId(stack.Name, update)
	})
	return
}

func (d *Dao) StackDelete(name string) (err error) {
	d.do(func(db *database) {
		err = db.C("stack").RemoveId(name)
	})
	return
}

// StackMigrate migrates stacks from old archive collection.
func (d *Dao) StackMigrate() {
	d.do(func(db *database) {
		logger := log.Get(app.Name)
		archiveColl := db.C("archive")

		// check collection is exists.
		if _, err := archiveColl.Indexes(); err != nil {
			return
		}

		archives := make([]*model.Archive, 0)
		err := archiveColl.Find(nil).All(&archives)
		if err != nil {
			logger.Warn("Failed to migrate archives: ", err)
			return
		}

		var errs []error
		stackColl := db.C("stack")
		for _, archive := range archives {
			stack := &model.Stack{
				Name:      archive.Name,
				Content:   archive.Content,
				CreatedBy: archive.CreatedBy,
				CreatedAt: archive.CreatedAt,
				UpdatedBy: archive.UpdatedBy,
				UpdatedAt: archive.UpdatedAt,
			}
			err = stackColl.Insert(stack)
			if err == nil || mgo.IsDup(err) {
				archiveColl.RemoveId(archive.ID)
			} else {
				logger.Warnf("Failed to migrate archive '%s': %v", archive.Name, err)
				errs = append(errs, err)
			}
		}

		// drop archive collection
		if len(errs) == 0 {
			err = archiveColl.DropCollection()
			if err != nil {
				logger.Warn("Failed to drop archive collection: ", err)
				return
			}
		}
	})
	return
}

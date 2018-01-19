package mongo

import (
	"time"

	"github.com/cuigh/swirl/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (d *Dao) TemplateList(args *model.TemplateListArgs) (tpls []*model.Template, count int, err error) {
	d.do(func(db *database) {
		filter := bson.M{}
		if args.Name != "" {
			filter["name"] = args.Name
		}

		q := db.C("template").Find(filter)
		count, err = q.Count()
		if err != nil {
			return
		}

		tpls = []*model.Template{}
		err = q.Skip(args.PageSize * (args.PageIndex - 1)).Limit(args.PageSize).All(&tpls)
	})
	return
}

func (d *Dao) TemplateCreate(tpl *model.Template) (err error) {
	tpl.CreatedAt = time.Now()
	tpl.UpdatedAt = tpl.CreatedAt

	d.do(func(db *database) {
		err = db.C("template").Insert(tpl)
	})
	return
}

func (d *Dao) TemplateGet(id string) (tpl *model.Template, err error) {
	d.do(func(db *database) {
		tpl = &model.Template{}
		err = db.C("template").FindId(id).One(tpl)
		if err == mgo.ErrNotFound {
			tpl, err = nil, nil
		} else if err != nil {
			tpl = nil
		}
	})
	return
}

func (d *Dao) TemplateUpdate(tpl *model.Template) (err error) {
	d.do(func(db *database) {
		update := bson.M{
			"$set": bson.M{
				"name":       tpl.Name,
				"content":    tpl.Content,
				"updated_by": tpl.UpdatedBy,
				"updated_at": tpl.UpdatedAt,
			},
		}
		err = db.C("template").UpdateId(tpl.ID, update)
	})
	return
}

func (d *Dao) TemplateDelete(id string) (err error) {
	d.do(func(db *database) {
		err = db.C("template").RemoveId(id)
	})
	return
}

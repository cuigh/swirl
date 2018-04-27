package bolt

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

func (d *Dao) TemplateList(args *model.TemplateListArgs) (tpls []*model.Template, count int, err error) {
	err = d.each("template", func(v Value) error {
		t := &model.Template{}
		err = v.Unmarshal(t)
		if err != nil {
			return err
		}

		if matchAny(args.Name, t.Name) {
			tpls = append(tpls, t)
		}
		return nil
	})
	if err == nil {
		count = len(tpls)
		start, end := misc.Page(count, args.PageIndex, args.PageSize)
		tpls = tpls[start:end]
	}
	return
}

func (d *Dao) TemplateCreate(tpl *model.Template) (err error) {
	tpl.CreatedAt = time.Now()
	tpl.UpdatedAt = tpl.CreatedAt
	return d.update("template", tpl.ID, tpl)
}

func (d *Dao) TemplateGet(id string) (tpl *model.Template, err error) {
	var v Value
	v, err = d.get("template", id)
	if err == nil {
		if v != nil {
			tpl = &model.Template{}
			err = v.Unmarshal(tpl)
		}
	}
	return
}

func (d *Dao) TemplateUpdate(tpl *model.Template) (err error) {
	return d.batch("template", func(b *bolt.Bucket) error {
		data := b.Get([]byte(tpl.ID))
		if data == nil {
			return errors.New("template not found: " + tpl.ID)
		}

		t := &model.Template{}
		err = json.Unmarshal(data, t)
		if err != nil {
			return err
		}

		t.Name = tpl.Name
		t.Content = tpl.Content
		t.UpdatedBy = tpl.UpdatedBy
		t.UpdatedAt = time.Now()
		data, err = json.Marshal(t)
		if err != nil {
			return err
		}

		return b.Put([]byte(tpl.ID), data)
	})
}

func (d *Dao) TemplateDelete(id string) (err error) {
	return d.delete("template", id)
}

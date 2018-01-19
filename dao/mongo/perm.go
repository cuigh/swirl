package mongo

import (
	"github.com/cuigh/swirl/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (d *Dao) PermGet(resType, resID string) (p *model.Perm, err error) {
	d.do(func(db *database) {
		p = &model.Perm{}
		q := bson.M{
			"res_type": resType,
			"res_id":   resID,
		}
		err = db.C("perm").Find(q).One(p)
		if err == mgo.ErrNotFound {
			p, err = nil, nil
		} else if err != nil {
			p = nil
		}
	})
	return
}

func (d *Dao) PermUpdate(perm *model.Perm) (err error) {
	d.do(func(db *database) {
		q := bson.M{
			"res_type": perm.ResType,
			"res_id":   perm.ResID,
		}
		update := bson.M{
			"$set": bson.M{
				"scope": perm.Scope,
				"roles": perm.Roles,
				"users": perm.Users,
			},
		}
		_, err = db.C("perm").Upsert(q, update)
	})
	return
}

func (d *Dao) PermDelete(resType, resID string) (err error) {
	d.do(func(db *database) {
		q := bson.M{
			"res_type": resType,
			"res_id":   resID,
		}
		err = db.C("perm").Remove(q)
	})
	return
}

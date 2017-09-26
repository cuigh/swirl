package mongo

import (
	"github.com/cuigh/swirl/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (d *Dao) RoleList() (roles []*model.Role, err error) {
	d.do(func(db *database) {
		roles = []*model.Role{}
		err = db.C("role").Find(nil).All(&roles)
	})
	return
}

func (d *Dao) RoleCreate(role *model.Role) (err error) {
	d.do(func(db *database) {
		err = db.C("role").Insert(role)
	})
	return
}

func (d *Dao) RoleGet(id string) (role *model.Role, err error) {
	d.do(func(db *database) {
		role = &model.Role{}
		err = db.C("role").FindId(id).One(role)
		if err == mgo.ErrNotFound {
			role, err = nil, nil
		} else if err != nil {
			role = nil
		}
	})
	return
}

func (d *Dao) RoleUpdate(role *model.Role) (err error) {
	d.do(func(db *database) {
		update := bson.M{
			"$set": bson.M{
				"name":       role.Name,
				"desc":       role.Description,
				"perms":      role.Perms,
				"updated_at": role.UpdatedAt,
			},
		}
		err = db.C("role").UpdateId(role.ID, update)
	})
	return
}

func (d *Dao) RoleDelete(id string) (err error) {
	d.do(func(db *database) {
		err = db.C("role").RemoveId(id)
		if err == nil {
			update := bson.M{
				"$pull": bson.M{"roles": id},
			}
			_, err = db.C("user").UpdateAll(nil, update)
		}
	})
	return
}

package mongo

import (
	"time"

	"github.com/cuigh/swirl/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (d *Dao) UserCount() (count int, err error) {
	d.do(func(db *database) {
		count, err = db.C("user").Count()
	})
	return
}

func (d *Dao) UserCreate(user *model.User) (err error) {
	d.do(func(db *database) {
		err = db.C("user").Insert(user)
	})
	return
}

func (d *Dao) UserUpdate(user *model.User) (err error) {
	d.do(func(db *database) {
		update := bson.M{
			"$set": bson.M{
				"name":       user.Name,
				"email":      user.Email,
				"admin":      user.Admin,
				"type":       user.Type,
				"roles":      user.Roles,
				"updated_at": user.UpdatedAt,
			},
		}
		err = db.C("user").UpdateId(user.ID, update)
	})
	return
}

func (d *Dao) UserBlock(id string, blocked bool) (err error) {
	d.do(func(db *database) {
		var status model.UserStatus
		if blocked {
			status = model.UserStatusBlocked
		} else {
			status = model.UserStatusActive
		}
		update := bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": time.Now(),
			},
		}
		err = db.C("user").UpdateId(id, update)
	})
	return
}

func (d *Dao) UserDelete(id string) (err error) {
	d.do(func(db *database) {
		err = db.C("user").RemoveId(id)
	})
	return
}

func (d *Dao) UserList(args *model.UserListArgs) (users []*model.User, count int, err error) {
	d.do(func(db *database) {
		query := bson.M{}
		if args.Query != "" {
			query["$or"] = []bson.M{
				{"login_name": args.Query},
				{"name": args.Query},
				{"email": args.Query},
			}
		}
		switch args.Filter {
		case "admins":
			query["admin"] = true
		case "active":
			query["status"] = 1
		case "blocked":
			query["status"] = 0
		}

		q := db.C("user").Find(query)
		count, err = q.Count()
		if err != nil {
			return
		}

		users = []*model.User{}
		err = q.Skip(args.PageSize * (args.PageIndex - 1)).Limit(args.PageSize).All(&users)
	})
	return
}

func (d *Dao) UserGetByID(id string) (user *model.User, err error) {
	d.do(func(db *database) {
		user = &model.User{}
		err = db.C("user").FindId(id).One(user)
		if err == mgo.ErrNotFound {
			user, err = nil, nil
		}
	})
	return
}

func (d *Dao) UserGetByName(loginName string) (user *model.User, err error) {
	d.do(func(db *database) {
		user = &model.User{}
		err = db.C("user").Find(bson.M{"login_name": loginName}).One(user)
		if err == mgo.ErrNotFound {
			user, err = nil, nil
		}
	})
	return
}

func (d *Dao) ProfileUpdateInfo(user *model.User) (err error) {
	d.do(func(db *database) {
		update := bson.M{
			"$set": bson.M{
				"name":       user.Name,
				"email":      user.Email,
				"updated_at": time.Now(),
			},
		}
		err = db.C("user").UpdateId(user.ID, update)
	})
	return
}

func (d *Dao) ProfileUpdatePassword(id, pwd, salt string) (err error) {
	d.do(func(db *database) {
		update := bson.M{
			"$set": bson.M{
				"password":   pwd,
				"salt":       salt,
				"updated_at": time.Now(),
			},
		}
		err = db.C("user").UpdateId(id, update)
	})
	return
}

func (d *Dao) SessionGet(token string) (session *model.Session, err error) {
	d.do(func(db *database) {
		session = &model.Session{}
		err = db.C("session").Find(bson.M{"token": token}).One(session)
		if err == mgo.ErrNotFound {
			session, err = nil, nil
		}
	})
	return
}

func (d *Dao) SessionUpdate(session *model.Session) (err error) {
	d.do(func(db *database) {
		_, err = db.C("session").UpsertId(session.UserID, session)
	})
	return
}

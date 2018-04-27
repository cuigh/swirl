package bolt

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

func (d *Dao) UserCount() (count int, err error) {
	return d.count("user")
}

func (d *Dao) UserCreate(user *model.User) (err error) {
	return d.update("user", user.ID, user)
}

func (d *Dao) UserUpdate(user *model.User) (err error) {
	return d.userUpdate(user.ID, func(u *model.User) {
		u.Name = user.Name
		u.Email = user.Email
		u.Admin = user.Admin
		u.Type = user.Type
		u.Roles = user.Roles
	})
}

func (d *Dao) UserBlock(id string, blocked bool) (err error) {
	return d.userUpdate(id, func(u *model.User) {
		if blocked {
			u.Status = model.UserStatusBlocked
		} else {
			u.Status = model.UserStatusActive
		}
	})
}

func (d *Dao) UserDelete(id string) (err error) {
	return d.delete("user", id)
}

func (d *Dao) UserList(args *model.UserListArgs) (users []*model.User, count int, err error) {
	err = d.each("user", func(v Value) error {
		user := &model.User{}
		err = v.Unmarshal(user)
		if err != nil {
			return err
		}

		match := true
		if args.Query != "" {
			match = matchAny(args.Query, user.LoginName, user.Name, user.Email)
		}
		if match {
			switch args.Filter {
			case "admins":
				match = user.Admin
			case "active":
				match = user.Status == model.UserStatusActive
			case "blocked":
				match = user.Status == model.UserStatusBlocked
			}
		}

		if match {
			users = append(users, user)
		}
		return nil
	})
	if err == nil {
		count = len(users)
		start, end := misc.Page(count, args.PageIndex, args.PageSize)
		users = users[start:end]
	}
	return
}

func (d *Dao) UserGetByID(id string) (user *model.User, err error) {
	var v Value
	v, err = d.get("user", id)
	if err == nil {
		if v != nil {
			user = &model.User{}
			err = v.Unmarshal(user)
		}
	}
	return
}

func (d *Dao) UserGetByName(loginName string) (user *model.User, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("user"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			u := &model.User{}
			err = json.Unmarshal(v, u)
			if err != nil {
				return err
			}
			if u.LoginName == loginName {
				user = u
				return nil
			}
		}
		return nil
	})
	return
}

func (d *Dao) ProfileUpdateInfo(user *model.User) (err error) {
	return d.userUpdate(user.ID, func(u *model.User) {
		u.Name = user.Name
		u.Email = user.Email
	})
}

func (d *Dao) ProfileUpdatePassword(id, pwd, salt string) (err error) {
	return d.userUpdate(id, func(u *model.User) {
		u.Password = pwd
		u.Salt = salt
	})
}

func (d *Dao) userUpdate(id string, decorator func(u *model.User)) (err error) {
	return d.batch("user", func(b *bolt.Bucket) error {
		data := b.Get([]byte(id))
		if data == nil {
			return errors.New("user not found: " + id)
		}

		u := &model.User{}
		err = json.Unmarshal(data, u)
		if err != nil {
			return err
		}

		decorator(u)
		u.UpdatedAt = time.Now()
		data, err = json.Marshal(u)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), data)
	})
}

func (d *Dao) SessionGet(token string) (session *model.Session, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("session"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			s := &model.Session{}
			err = json.Unmarshal(v, s)
			if err != nil {
				return err
			}
			if s.Token == token {
				session = s
				return nil
			}
		}
		return nil
	})
	return
}

func (d *Dao) SessionUpdate(session *model.Session) (err error) {
	return d.update("session", session.UserID, session)
}

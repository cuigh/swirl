package bolt

import (
	"context"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

func (d *Dao) UserCount(ctx context.Context) (count int, err error) {
	return d.count("user")
}

func (d *Dao) UserCreate(ctx context.Context, user *model.User) (err error) {
	return d.update("user", user.ID, user)
}

func (d *Dao) UserUpdate(ctx context.Context, user *model.User) (err error) {
	return d.userUpdate(user.ID, func(u *model.User) {
		u.Name = user.Name
		u.LoginName = user.LoginName
		u.Email = user.Email
		u.Admin = user.Admin
		u.Type = user.Type
		u.Roles = user.Roles
	})
}

func (d *Dao) UserSetStatus(ctx context.Context, id string, status int32) (err error) {
	return d.userUpdate(id, func(u *model.User) {
		u.Status = status
	})
}

func (d *Dao) UserDelete(ctx context.Context, id string) (err error) {
	return d.delete("user", id)
}

func (d *Dao) UserList(ctx context.Context, args *model.UserSearchArgs) (users []*model.User, count int, err error) {
	err = d.each("user", func(v Value) error {
		user := &model.User{}
		err = v.Unmarshal(user)
		if err != nil {
			return err
		}

		match := true
		if args.Name != "" {
			match = matchAny(args.Name, user.Name)
		}
		if match && args.LoginName != "" {
			match = matchAny(args.LoginName, user.LoginName)
		}
		if match && args.Admin {
			match = user.Admin
		}
		if match && args.Status >= 0 {
			match = user.Status == args.Status
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

func (d *Dao) UserGetByID(ctx context.Context, id string) (user *model.User, err error) {
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

func (d *Dao) UserGetByName(ctx context.Context, loginName string) (user *model.User, err error) {
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

func (d *Dao) UserModifyProfile(ctx context.Context, user *model.User) (err error) {
	return d.userUpdate(user.ID, func(u *model.User) {
		u.Name = user.Name
		u.LoginName = user.LoginName
		u.Email = user.Email
	})
}

func (d *Dao) UserModifyPassword(ctx context.Context, id, pwd, salt string) (err error) {
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

func (d *Dao) SessionGet(ctx context.Context, token string) (session *model.Session, err error) {
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

func (d *Dao) SessionUpdate(ctx context.Context, session *model.Session) (err error) {
	return d.update("session", session.UserID, session)
}

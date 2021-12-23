package bolt

import (
	"context"

	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/misc"
)

const User = "user"

func (d *Dao) UserCount(ctx context.Context) (count int, err error) {
	return d.count(User)
}

func (d *Dao) UserCreate(ctx context.Context, user *dao.User) (err error) {
	return d.replace(User, user.ID, user)
}

func (d *Dao) UserUpdate(ctx context.Context, user *dao.User) (err error) {
	old := &dao.User{}
	return d.update(User, user.ID, old, func() interface{} {
		old.Name = user.Name
		old.LoginName = user.LoginName
		old.Email = user.Email
		old.Admin = user.Admin
		old.Type = user.Type
		old.Roles = user.Roles
		old.UpdatedAt = user.UpdatedAt
		old.UpdatedBy = user.UpdatedBy
		return old
	})
}

func (d *Dao) UserUpdateStatus(ctx context.Context, user *dao.User) (err error) {
	old := &dao.User{}
	return d.update(User, user.ID, old, func() interface{} {
		old.Status = user.Status
		old.UpdatedAt = user.UpdatedAt
		old.UpdatedBy = user.UpdatedBy
		return old
	})
}

func (d *Dao) UserDelete(ctx context.Context, id string) (err error) {
	return d.delete(User, id)
}

func (d *Dao) UserSearch(ctx context.Context, args *dao.UserSearchArgs) (users []*dao.User, count int, err error) {
	err = d.each(User, func(v []byte) error {
		user := &dao.User{}
		err = decode(v, user)
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

func (d *Dao) UserGet(ctx context.Context, id string) (user *dao.User, err error) {
	user = &dao.User{}
	err = d.get(User, id, user)
	if err == ErrNoRecords {
		return nil, nil
	} else if err != nil {
		user = nil
	}
	return
}

func (d *Dao) UserGetByName(ctx context.Context, loginName string) (user *dao.User, err error) {
	u := &dao.User{}
	found, err := d.find(User, u, func() bool { return u.LoginName == loginName })
	if found {
		return u, nil
	}
	return nil, err
}

func (d *Dao) UserUpdateProfile(ctx context.Context, user *dao.User) (err error) {
	old := &dao.User{}
	return d.update(User, user.ID, old, func() interface{} {
		old.Name = user.Name
		old.LoginName = user.LoginName
		old.Email = user.Email
		old.UpdatedAt = user.UpdatedAt
		old.UpdatedBy = user.UpdatedBy
		return old
	})
}

func (d *Dao) UserUpdatePassword(ctx context.Context, user *dao.User) (err error) {
	old := &dao.User{}
	return d.update(User, user.ID, old, func() interface{} {
		old.Password = user.Password
		old.Salt = user.Salt
		old.UpdatedAt = user.UpdatedAt
		old.UpdatedBy = user.UpdatedBy
		return old
	})
}

package mongo

import (
	"context"

	"github.com/cuigh/swirl/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	User    = "user"
	Session = "session"
)

func (d *Dao) UserCount(ctx context.Context) (int, error) {
	count, err := d.db.Collection(User).CountDocuments(ctx, bson.M{})
	return int(count), err
}

func (d *Dao) UserCreate(ctx context.Context, user *model.User) (err error) {
	return d.create(ctx, User, user)
}

func (d *Dao) UserUpdate(ctx context.Context, user *model.User) (err error) {
	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"login_name": user.LoginName,
			"email":      user.Email,
			"admin":      user.Admin,
			"type":       user.Type,
			"roles":      user.Roles,
			"updated_at": user.UpdatedAt,
			"updated_by": user.UpdatedBy,
		},
	}
	return d.update(ctx, User, user.ID, update)
}

func (d *Dao) UserUpdateStatus(ctx context.Context, user *model.User) (err error) {
	update := bson.M{
		"$set": bson.M{
			"status":     user.Status,
			"updated_at": user.UpdatedAt,
			"updated_by": user.UpdatedBy,
		},
	}
	return d.update(ctx, User, user.ID, update)
}

func (d *Dao) UserDelete(ctx context.Context, id string) (err error) {
	return d.delete(ctx, User, id)
}

func (d *Dao) UserSearch(ctx context.Context, args *model.UserSearchArgs) (users []*model.User, count int, err error) {
	filter := bson.M{}
	if args.Name != "" {
		filter["name"] = args.Name
	}
	if args.LoginName != "" {
		filter["login_name"] = args.LoginName
	}
	if args.Admin {
		filter["admin"] = true
	}
	if args.Status >= 0 {
		filter["status"] = args.Status
	}
	opts := searchOptions{filter: filter, pageIndex: args.PageIndex, pageSize: args.PageSize}
	users = []*model.User{}
	count, err = d.search(ctx, User, opts, &users)
	return
}

func (d *Dao) UserGet(ctx context.Context, id string) (user *model.User, err error) {
	user = &model.User{}
	found, err := d.find(ctx, User, id, user)
	if !found {
		return nil, err
	}
	return
}

func (d *Dao) UserGetByName(ctx context.Context, loginName string) (user *model.User, err error) {
	user = &model.User{}
	err = d.db.Collection(User).FindOne(ctx, bson.M{"login_name": loginName}).Decode(user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return
}

func (d *Dao) UserUpdateProfile(ctx context.Context, user *model.User) (err error) {
	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"login_name": user.LoginName,
			"email":      user.Email,
			"updated_at": user.UpdatedAt,
			"updated_by": user.UpdatedBy,
		},
	}
	return d.update(ctx, User, user.ID, update)
}

func (d *Dao) UserUpdatePassword(ctx context.Context, user *model.User) (err error) {
	update := bson.M{
		"$set": bson.M{
			"password":   user.Password,
			"salt":       user.Salt,
			"updated_at": user.UpdatedAt,
			"updated_by": user.UpdatedBy,
		},
	}
	return d.update(ctx, User, user.ID, update)
}

func (d *Dao) SessionGet(ctx context.Context, token string) (session *model.Session, err error) {
	session = &model.Session{}
	err = d.db.Collection(Session).FindOne(ctx, bson.M{"token": token}).Decode(session)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return
}

func (d *Dao) SessionUpdate(ctx context.Context, session *model.Session) (err error) {
	_, err = d.db.Collection(Session).UpdateByID(ctx, session.UserID, session, options.Update().SetUpsert(true))
	return
}

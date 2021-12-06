package mongo

import (
	"context"

	"github.com/cuigh/swirl/model"
	"go.mongodb.org/mongo-driver/bson"
)

const Role = "role"

func (d *Dao) RoleList(ctx context.Context, name string) (roles []*model.Role, err error) {
	filter := bson.M{}
	if name != "" {
		filter["name"] = name
	}
	roles = []*model.Role{}
	err = d.fetch(ctx, Role, filter, &roles)
	return
}

func (d *Dao) RoleCreate(ctx context.Context, role *model.Role) (err error) {
	return d.create(ctx, Role, role)
}

func (d *Dao) RoleGet(ctx context.Context, id string) (role *model.Role, err error) {
	role = &model.Role{}
	found, err := d.find(ctx, Role, id, role)
	if !found {
		return nil, err
	}
	return
}

func (d *Dao) RoleUpdate(ctx context.Context, role *model.Role) (err error) {
	update := bson.M{
		"$set": bson.M{
			"name":       role.Name,
			"desc":       role.Description,
			"perms":      role.Perms,
			"updated_at": role.UpdatedAt,
		},
	}
	return d.update(ctx, Role, role.ID, update)
}

func (d *Dao) RoleDelete(ctx context.Context, id string) (err error) {
	err = d.delete(ctx, Role, id)
	if err == nil {
		update := bson.M{
			"$pull": bson.M{"roles": id},
		}
		_, err = d.db.Collection(User).UpdateMany(ctx, bson.M{}, update)
	}
	return
}

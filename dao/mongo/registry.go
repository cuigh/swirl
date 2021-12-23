package mongo

import (
	"context"

	"github.com/cuigh/swirl/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const Registry = "registry"

func (d *Dao) RegistryCreate(ctx context.Context, registry *dao.Registry) (err error) {
	return d.create(ctx, Registry, registry)
}

func (d *Dao) RegistryUpdate(ctx context.Context, registry *dao.Registry) (err error) {
	update := bson.M{
		"name":       registry.Name,
		"url":        registry.URL,
		"username":   registry.Username,
		"updated_at": registry.UpdatedAt,
		"updated_by": registry.UpdatedBy,
	}
	if registry.Password != "" {
		update["password"] = registry.Password
	}
	return d.update(ctx, Registry, registry.ID, bson.M{"$set": update})
}

func (d *Dao) RegistryGetAll(ctx context.Context) (registries []*dao.Registry, err error) {
	registries = []*dao.Registry{}
	err = d.fetch(ctx, Registry, bson.M{}, &registries)
	return
}

func (d *Dao) RegistryGet(ctx context.Context, id string) (registry *dao.Registry, err error) {
	registry = &dao.Registry{}
	found, err := d.find(ctx, Registry, id, registry)
	if !found {
		return nil, err
	}
	return
}

func (d *Dao) RegistryGetByURL(ctx context.Context, url string) (registry *dao.Registry, err error) {
	registry = &dao.Registry{}
	err = d.db.Collection(Registry).FindOne(ctx, bson.M{"url": url}).Decode(registry)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return
}

func (d *Dao) RegistryDelete(ctx context.Context, id string) (err error) {
	return d.delete(ctx, Registry, id)
}

// Package collections: used for handleEvent to communicate with the database
package collections

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"hunt/db"
	"hunt/structs"
)

type LocCollection struct {
	Collection *db.Repository[structs.HuntUser]
}

var LocationCollection = newLoctionCollection()

func newLoctionCollection() *LocCollection {
	repo, err := db.NewRepository[structs.HuntUser](db.HuntArtakClient, "location", nil)
	if err != nil {
		panic(err)
	}
	collection := &repo
	return &LocCollection{
		Collection: collection,
	}
}

func (lc *LocCollection) FindById(ctx context.Context, id string) (structs.HuntUser, error) {
	filter := db.CreateFilter("uid", id)
	result, err := lc.Collection.FindOne(ctx, filter)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (lc *LocCollection) UpdateUser(ctx context.Context, usr structs.HuntUser) error {
	filter := db.CreateFilter("uid", usr.Uid)
	_, err := lc.Collection.Upsert(ctx, filter, usr)
	if err != nil {
		return err
	}
	return nil
}

func (lc *LocCollection) AddUser(ctx context.Context, usr structs.HuntUser) error {
	_, err := lc.Collection.InsertOne(ctx, usr)
	if err != nil {
		return err
	}
	return nil
}

func (lc *LocCollection) SetInactive(ctx context.Context, uid string) error {
	filter := db.CreateFilter("uid", uid)
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "active",
					Value: false,
				},
			},
		},
	}
	err := lc.Collection.UpdateFields(ctx, filter, update)
	if err.Err() != nil {
		return err.Err()
	}
	return nil
}

func (lc *LocCollection) SetCallsign(ctx context.Context, uid, callsign string) error {
	filter := db.CreateFilter("uid", uid)
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "callsign",
					Value: callsign,
				},
			},
		},
	}
	err := lc.Collection.UpdateFields(ctx, filter, update)
	if err.Err() != nil {
		return err.Err()
	}
	return nil
}

func (lc *LocCollection) SetActive(ctx context.Context, uid string) error {
	filter := db.CreateFilter("uid", uid)
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "active",
					Value: true,
				},
			},
		},
	}
	err := lc.Collection.UpdateFields(ctx, filter, update)
	if err.Err() != nil {
		return err.Err()
	}
	return nil
}

func (lc *LocCollection) UpdateTarget(ctx context.Context, uid, targetUid, targetName string) error {
	filter := db.CreateFilter("uid", uid)
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "targetuid",
					Value: targetUid,
				},
				{
					Key:   "targetname",
					Value: targetName,
				},
			},
		},
	}
	err := lc.Collection.UpdateFields(ctx, filter, update)
	if err.Err() != nil {
		return err.Err()
	}
	return nil
}

func (lc *LocCollection) FindAll(ctx context.Context) ([]structs.HuntUser, error) {
	result, err := lc.Collection.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (lc *LocCollection) FindAllActive(ctx context.Context) ([]structs.HuntUser, error) {
	filter := bson.M{
		"active": true,
	}
	result, err := lc.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (lc *LocCollection) ReadManyLocations(ctx context.Context, uids, partitionKeys []string) ([]structs.HuntUser, error) {
	filter := bson.M{
		"$and": bson.A{
			bson.M{"uid": bson.M{"$in": uids}},
			bson.M{"$or": bson.A{
				bson.M{"partition_keys": bson.M{"$exists": false}},
				bson.M{"partition_keys": bson.M{"in": partitionKeys}},
			}},
		},
	}
	result, err := lc.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

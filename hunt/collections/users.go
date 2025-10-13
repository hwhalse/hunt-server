// Package collections: used for handleEvent to communicate with the database
package collections

import (
	"context"
	"errors"
	"hunt/db"
	"hunt/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersCollection struct {
	Collection *db.Repository[models.HuntUser]
}

var UsersCollectionManager = newUsersCollectionManager()

func newUsersCollectionManager() *UsersCollection {
	repo, err := db.NewRepository[models.HuntUser](db.HuntArtakClient, "location", nil)
	if err != nil {
		panic(err)
	}
	collection := &repo
	return &UsersCollection{
		Collection: collection,
	}
}

func (lc *UsersCollection) FindById(ctx context.Context, id string) (models.HuntUser, error) {
	filter := db.CreateFilter("uid", id)
	result, err := lc.Collection.FindOne(ctx, filter)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (lc *UsersCollection) UpdateLocation(ctx context.Context, location models.Location, uid string) error {
	filter := bson.M{"uid": uid}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "location", Value: location},
			},
		},
	}

	result := lc.Collection.UpdateFields(ctx, filter, update)

	// Check for internal error (e.g., invalid query)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Fatalf("CRITICAL: No document found for uid: %s", uid)
		}
	}

	var updated models.HuntUser
	if err := result.Decode(&updated); err != nil {
		log.Fatalf("Failed to decode updated user: %v\n", err)
	}

	return nil
}

func (lc *UsersCollection) UpdateUser(ctx context.Context, usr models.HuntUser) error {
	filter := db.CreateFilter("uid", usr.Uid)
	_, err := lc.Collection.Upsert(ctx, filter, usr)
	if err != nil {
		return err
	}
	return nil
}

func (lc *UsersCollection) AddUser(ctx context.Context, usr models.HuntUser) error {
	_, err := lc.Collection.InsertOne(ctx, usr)
	if err != nil {
		return err
	}
	return nil
}

func (lc *UsersCollection) SetInactive(ctx context.Context, uid string) error {
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

func (lc *UsersCollection) SetCallsign(ctx context.Context, uid, callsign string) error {
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

func (lc *UsersCollection) SetActive(ctx context.Context, uid string) error {
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

func (lc *UsersCollection) UpdateTarget(ctx context.Context, uid, targetUid, targetName string) error {
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

func (lc *UsersCollection) FindAll(ctx context.Context) ([]models.HuntUser, error) {
	result, err := lc.Collection.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (lc *UsersCollection) FindAllActive(ctx context.Context) ([]models.HuntUser, error) {
	filter := bson.M{
		"active": true,
	}
	result, err := lc.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (lc *UsersCollection) ReadManyLocations(ctx context.Context, uids, partitionKeys []string) ([]models.HuntUser, error) {
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

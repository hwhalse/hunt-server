// Package collections: used for handleEvent to communicate with the database
package collections

import (
	"context"
	"hunt/db"
	"hunt/logging"
	"hunt/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

var CommandNodeCollectionManager = newCommandNodeCollectionManager()

func newCommandNodeCollectionManager() *CommandNodeCollection {
	repo, err := db.NewRepository[models.CommandNode](db.HuntArtakClient, "command_nodes", nil)
	if err != nil {
		panic(err)
	}
	collection := CommandNodeCollection{
		Collection: &repo,
		Logger:     logging.NewLogger(),
	}
	return &collection
}

type CommandNodeCollection struct {
	Collection *db.Repository[models.CommandNode]
	Logger     zerolog.Logger
}

func (cmn *CommandNodeCollection) FindAll(ctx context.Context) ([]models.CommandNode, error) {
	res, err := cmn.Collection.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cmn *CommandNodeCollection) FindById(ctx context.Context, id string) (models.CommandNode, error) {
	filter := db.CreateFilter("uid", id)
	res, err := cmn.Collection.FindOne(ctx, filter)
	if err != nil {
		return models.CommandNode{}, err
	}
	return res, err
}

func (cmn *CommandNodeCollection) SetStatus(ctx context.Context, id string, status int) error {
	filter := db.CreateFilter("uid", id)
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "status",
					Value: status,
				},
			},
		},
	}
	res := cmn.Collection.UpdateFields(ctx, filter, update)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (cmn *CommandNodeCollection) InsertNode(ctx context.Context, node models.CommandNode) error {
	if node.Uid == "" {
		node.Uid = uuid.NewString()
	}
	_, err := cmn.Collection.InsertOne(ctx, node)
	if err != nil {
		return err
	}
	return nil
}

func (cmn *CommandNodeCollection) Delete(ctx context.Context, filter bson.M) error {
	err := cmn.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (cmn *CommandNodeCollection) FindOneAndUpdate(ctx context.Context, filter bson.M, update models.CommandNode) (models.CommandNode, error) {
	res, err := cmn.Collection.Upsert(ctx, filter, update)
	if err != nil {
		cmn.Logger.Error().Err(err).Msg("find one and update error")
		return res, err
	}
	return res, nil
}

func (cmn *CommandNodeCollection) ReadManyNodes(ctx context.Context, uids []string, partitionKeys []string) ([]models.CommandNode, error) {
	filter := bson.M{
		"$and": bson.A{
			bson.M{"uid": bson.M{"$in": uids}},
			bson.M{"$or": bson.A{
				bson.M{"partition_keys": bson.M{"$exists": false}},
				bson.M{"partition_keys": bson.M{"in": partitionKeys}},
			}},
		},
	}
	var emptyNodes []models.CommandNode
	result, err := cmn.Collection.Find(ctx, filter)
	if err != nil {
		return emptyNodes, nil
	}
	return result, nil
}

// Package collections: used for handleEvent to communicate with the database
package collections

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"hunt/db"
	"hunt/logging"
	"hunt/structs"
)

var CommandNodeCollection = newCmdNodeCollection()

func newCmdNodeCollection() *CmdNodeCollection {
	repo, err := db.NewRepository[structs.CommandNode](db.HuntArtakClient, "command_nodes", nil)
	if err != nil {
		panic(err)
	}
	collection := CmdNodeCollection{
		Collection: &repo,
		Logger:     logging.NewLogger(),
	}
	return &collection
}

type CmdNodeCollection struct {
	Collection *db.Repository[structs.CommandNode]
	Logger     zerolog.Logger
}

func (cmn *CmdNodeCollection) FindAll(ctx context.Context) ([]structs.CommandNode, error) {
	res, err := cmn.Collection.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cmn *CmdNodeCollection) FindById(ctx context.Context, id string) (structs.CommandNode, error) {
	filter := db.CreateFilter("uid", id)
	res, err := cmn.Collection.FindOne(ctx, filter)
	if err != nil {
		return structs.CommandNode{}, err
	}
	return res, err
}

func (cmn *CmdNodeCollection) SetStatus(ctx context.Context, id string, status int) error {
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

func (cmn *CmdNodeCollection) InsertNode(ctx context.Context, node structs.CommandNode) error {
	if node.Uid == "" {
		node.Uid = uuid.NewString()
	}
	_, err := cmn.Collection.InsertOne(ctx, node)
	if err != nil {
		return err
	}
	return nil
}

func (cmn *CmdNodeCollection) Delete(ctx context.Context, filter bson.M) error {
	err := cmn.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (cmn *CmdNodeCollection) FindOneAndUpdate(ctx context.Context, filter bson.M, update structs.CommandNode) (structs.CommandNode, error) {
	res, err := cmn.Collection.Upsert(ctx, filter, update)
	if err != nil {
		cmn.Logger.Error().Err(err).Msg("find one and update error")
		return res, err
	}
	return res, nil
}

func (cmn *CmdNodeCollection) ReadManyNodes(ctx context.Context, uids []string, partitionKeys []string) ([]structs.CommandNode, error) {
	filter := bson.M{
		"$and": bson.A{
			bson.M{"uid": bson.M{"$in": uids}},
			bson.M{"$or": bson.A{
				bson.M{"partition_keys": bson.M{"$exists": false}},
				bson.M{"partition_keys": bson.M{"in": partitionKeys}},
			}},
		},
	}
	var emptyNodes []structs.CommandNode
	result, err := cmn.Collection.Find(ctx, filter)
	if err != nil {
		return emptyNodes, nil
	}
	return result, nil
}

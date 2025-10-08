// Package collections: used for handleEvent to communicate with the database
package collections

import (
	"context"
	"hunt/constants"
	"hunt/db"
	"hunt/structs"
)

type GrpCollection struct {
	Collection *db.Repository[structs.Unit]
}

var GroupCollection = newGrpCollection()

func newGrpCollection() *GrpCollection {
	repo, err := db.NewRepository[structs.Unit](db.HuntArtakClient, "groups", nil)
	if err != nil {
		panic(err)
	}
	collection := &repo
	return &GrpCollection{
		Collection: collection,
	}
}

func (cmn *GrpCollection) FindAll(ctx context.Context) ([]structs.Unit, error) {
	result, err := cmn.Collection.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cmn *GrpCollection) InsertUnit(ctx context.Context, unit structs.Unit) error {
	inserted, err := cmn.Collection.InsertOne(ctx, unit)
	if err != nil {
		return err
	}
	if inserted.InsertedID == nil {
		return constants.ErrNoDocsInserted
	}
	return nil
}

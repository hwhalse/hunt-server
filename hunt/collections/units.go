// Package collections: used for handleEvent to communicate with the database
package collections

import (
	"context"
	"hunt/constants"
	"hunt/db"
	"hunt/models"
)

type UnitCollection struct {
	Collection *db.Repository[models.Unit]
}

var UnitCollectionManager = newUnitCollectionManager()

func newUnitCollectionManager() *UnitCollection {
	repo, err := db.NewRepository[models.Unit](db.HuntArtakClient, "units", nil)
	if err != nil {
		panic(err)
	}
	collection := &repo
	return &UnitCollection{
		Collection: collection,
	}
}

func (c *UnitCollection) FindAll(ctx context.Context) ([]models.Unit, error) {
	result, err := c.Collection.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *UnitCollection) InsertUnit(ctx context.Context, unit models.Unit) error {
	inserted, err := c.Collection.InsertOne(ctx, unit)
	if err != nil {
		return err
	}
	if inserted.InsertedID == nil {
		return constants.ErrNoDocsInserted
	}
	return nil
}

package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hunt/constants"
)

// Generic mongodb repository that provides base methods for each service's collections to use

type Repository[T any] struct {
	collection *mongo.Collection
}

func NewRepository[T any](client *Client, collectionName string, indexes []mongo.IndexModel) (Repository[T], error) {
	repo := Repository[T]{
		collection: client.Database().Collection(collectionName),
	}
	for _, index := range indexes {
		_, err := repo.collection.Indexes().CreateOne(context.TODO(), index)
		if err != nil {
			return repo, err
		}
	}
	return repo, nil
}

func (repo *Repository[T]) FindAll(ctx context.Context) ([]T, error) {
	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (repo *Repository[T]) Find(ctx context.Context, filter bson.M) ([]T, error) {
	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (repo *Repository[T]) FindOne(ctx context.Context, filter bson.M) (T, error) {
	var result T
	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	return result, err
}

// Upsert replaces entire document
func (repo *Repository[T]) Upsert(ctx context.Context, filter bson.M, document T) (T, error) {
	var result T
	opts := options.FindOneAndReplace().SetUpsert(true).SetReturnDocument(options.After)
	err := repo.collection.FindOneAndReplace(ctx, filter, document, opts).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return result, nil
	}
	return result, err
}

func (repo *Repository[T]) UpdateFields(ctx context.Context, filter bson.M, update bson.D) *mongo.SingleResult {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := repo.collection.FindOneAndUpdate(ctx, filter, update, opts)
	return result
}

func (repo *Repository[T]) DeleteMany(ctx context.Context, filter bson.M) error {
	result, err := repo.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return constants.ErrNoDocsDeleted
	}
	return nil
}

func (repo *Repository[T]) FindOneAndDelete(ctx context.Context, filter bson.M) error {
	result := repo.collection.FindOneAndDelete(ctx, filter)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (repo *Repository[T]) Aggregate(ctx context.Context, pipeline mongo.Pipeline) ([]T, error) {
	var results []T
	cursor, err := repo.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return results, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		return results, err
	}
	return results, nil
}

func (repo *Repository[T]) InsertOne(ctx context.Context, document T) (*mongo.InsertOneResult, error) {
	return repo.collection.InsertOne(ctx, document)
}

func (repo *Repository[T]) UpdateOne(ctx context.Context, filter bson.M, update bson.M, opts *options.UpdateOptions) error {
	_, err := repo.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (repo *Repository[T]) FindOneAndUpdate(ctx context.Context, filter bson.M, update bson.M, opts *options.FindOneAndUpdateOptions) (T, error) {
	var res T
	err := repo.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (repo *Repository[T]) Distinct(ctx context.Context, fieldName string, filter interface{}) ([]interface{}, error) {
	res, err := repo.collection.Distinct(ctx, fieldName, filter)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *Repository[T]) DeleteOne(ctx context.Context, filter bson.M) error {
	_, err := repo.collection.DeleteOne(ctx, filter)
	return err
}

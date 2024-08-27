package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBRepository[T any] struct {
	readWrapper  func(map[string]interface{}) (*T, error) // readWrapper reads a map and converts
	writeWrapper func(T) (map[string]interface{}, error)  // writeWrapper writes a value to a map

	collection *mongo.Collection
}

// Insert inserts a new value or update into the repository
func (r *MongoDBRepository[T]) Insert(v T) (*mongo.InsertOneResult, error) {
	data, err := r.writeWrapper(v)
	if err != nil {
		return nil, err
	}

	return r.collection.InsertOne(context.Background(), data)
}

// FindOne finds a value by its ID
func (r *MongoDBRepository[T]) FindOne(id string) (*T, error) {
	data := r.collection.FindOne(context.Background(), map[string]interface{}{"_id": id})
	if data.Err() != nil {
		return nil, data.Err()
	}

	var v map[string]interface{}
	if err := data.Decode(&v); err != nil {
		return nil, err
	}

	return r.readWrapper(v)
}

func NewMongoDB(
	readWrapper func(map[string]interface{}) (*T, error),
	writeWrapper func(T) (map[string]interface{}, error),
	collectionName string,
) *MongoDBRepository[T] {
	return &MongoDBRepository[T]{
		readWrapper:  readWrapper,
		writeWrapper: writeWrapper,

		collection: collection,
	}
}

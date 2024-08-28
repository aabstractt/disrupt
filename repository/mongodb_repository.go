package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBRepository[T any] struct {
	reader func(map[string]interface{}) (T, error) // readWrapper reads a map and converts
	writer func(T) (map[string]interface{}, error) // writeWrapper writes a value to a map
	errors func(error) (T, error)                  // errorWrapper returns an error

	collection *mongo.Collection
}

// Insert inserts a new value or update into the repository
func (r *MongoDBRepository[T]) Insert(v T) (*mongo.UpdateResult, error) {
	data, err := r.writer(v)
	if err != nil {
		return nil, err
	}

	return r.collection.UpdateOne(context.Background(), data)
}

func (r *MongoDBRepository[T]) Delete(id string) (*mongo.DeleteResult, error) {
	return r.collection.DeleteOne(context.Background(), map[string]interface{}{"_id": id})
}

// FindOne finds a value by its ID
func (r *MongoDBRepository[T]) FindOne(id string) (T, error) {
	data := r.collection.FindOne(context.Background(), map[string]interface{}{"_id": id})
	if data.Err() != nil {
		return r.errors(data.Err())
	}

	var v map[string]interface{}
	if err := data.Decode(&v); err != nil {
		return r.errors(err)
	}

	return r.reader(v)
}

func (r *MongoDBRepository[T]) FindMany(k, v string) ([]T, error) {
	result, err := r.collection.Find(context.Background(), map[string]interface{}{k: v})
	if err != nil {
		return nil, err
	}

	var values []T
	for result.Next(context.Background()) {
		var data map[string]interface{}
		if err := result.Decode(&data); err != nil {
			return nil, err
		}

		v, err := r.reader(data)
		if err != nil {
			return nil, err
		}

		values = append(values, v)
	}

	return values, nil
}

func (r *MongoDBRepository[T]) FindAll() ([]T, error) {
	result, err := r.collection.Find(context.Background(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var values []T
	for result.Next(context.Background()) {
		var data map[string]interface{}
		if err := result.Decode(&data); err != nil {
			return nil, err
		}

		v, err := r.reader(data)
		if err != nil {
			return nil, err
		}

		values = append(values, v)
	}

	return values, nil
}

func NewMongoDB[T any](
	readWrapper func(map[string]interface{}) (T, error),
	writeWrapper func(T) (map[string]interface{}, error),
	collectionName string,
) *MongoDBRepository[T] {
	return &MongoDBRepository[T]{
		reader: readWrapper,
		writer: writeWrapper,

		collection: collection,
	}
}

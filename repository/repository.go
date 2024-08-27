package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository[T any] interface {
	// Insert inserts a new value or update into the repository
	Insert(v T) (*mongo.InsertOneResult, error)

	// Delete deletes a value by its ID
	Delete(id string) (*mongo.DeleteResult, error)

	// FindOne finds a value by its ID
	FindOne(id string) (T, error)

	// FindMany finds values by a key and value
	FindMany(k, v string) ([]T, error)

	// FindAll finds all values in the repository
	FindAll() ([]T, error)
}

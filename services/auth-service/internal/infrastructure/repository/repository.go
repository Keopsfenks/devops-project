package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository[T any] struct {
	collection *mongo.Collection
}

func NewRepository[T any](collection *mongo.Collection) *Repository[T] {
	return &Repository[T]{collection: collection}
}

func (r *Repository[T]) Find(ctx context.Context, filter map[string]interface{}, skip, limit int) ([]*T, error) {

}

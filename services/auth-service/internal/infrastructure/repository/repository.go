package repository

import (
	"context"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository[T any] struct {
	collection *mongo.Collection
}

func NewRepository[T any](collection *mongo.Collection) *Repository[T] {
	return &Repository[T]{collection: collection}
}

func (r Repository[T]) Find(ctx context.Context, filter map[string]interface{}, skip, limit int) ([]*T, error) {
	if filter == nil {
		filter = make(map[string]interface{})
	}
	filter["is_deleted"] = false

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *Repository[T]) FindOne(ctx context.Context, filter map[string]interface{}) (*T, error) {
	if filter == nil {
		filter = make(map[string]interface{})
	}
	filter["is_deleted"] = false

	cursor := r.collection.FindOne(ctx, filter)
	var result T
	if err := cursor.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository[T]) InsertOne(ctx context.Context, entity *T) (*T, error) {
	result, err := r.collection.InsertOne(ctx, entity)

	if err != nil {
		log.Fatalf("failed to insert document: %v, error: %v", entity, err)
		return nil, err
	}

	if result.InsertedID == nil {
		log.Fatalf("failed to insert document: %v", entity)
		return nil, mongo.ErrNilDocument
	}
	log.Printf("inserted document: %v", entity)
	return entity, nil
}

func (r *Repository[T]) ReplaceOne(ctx context.Context, filter map[string]interface{}, entity *T) (*T, error) {
	now := time.Now()

	v := reflect.ValueOf(entity).Elem()
	field := v.FieldByName("UpdatedAt")
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(&now))
	}

	result, err := r.collection.ReplaceOne(ctx, filter, entity)
	if err != nil {
		log.Fatalf("failed to replace document: %v, error: %v", entity, err)
		return nil, err
	}

	if result.MatchedCount == 0 {
		log.Fatalf("no document matched the filter: %v", filter)
		return nil, mongo.ErrNoDocuments
	}

	log.Printf("replaced document: %v", entity)
	return entity, nil
}

func (r *Repository[T]) DeleteOne(ctx context.Context, filter map[string]interface{}) (bool, error) {
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalf("failed to delete document with filter: %v, error: %v", filter, err)
		return false, err
	}
	if result.DeletedCount == 0 {
		log.Fatalf("no document matched the filter: %v", filter)
		return false, mongo.ErrNoDocuments
	}
	log.Printf("deleted document with filter: %v", filter)
	return true, nil
}

func (r *Repository[T]) SoftDeleteOne(ctx context.Context, filter map[string]interface{}, entity *T) (*T, error) {
	now := time.Now()

	v := reflect.ValueOf(entity).Elem()

	field := v.FieldByName("DeletedAt")
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(&now))
	}

	result, err := r.collection.ReplaceOne(ctx, filter, entity)

	if err != nil {
		log.Fatalf("failed to soft delete document: %v, error: %v", entity, err)
		return nil, err
	}

	if result.MatchedCount == 0 {
		log.Fatalf("no document matched the filter: %v", filter)
		return nil, mongo.ErrNoDocuments
	}

	log.Printf("soft deleted document: %v", entity)
	return entity, nil
}

func (r *Repository[T]) Exists(ctx context.Context, filter map[string]interface{}) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatalf("failed to count documents with filter: %v, error: %v", filter, err)
		return false, err
	}
	return count > 0, nil
}

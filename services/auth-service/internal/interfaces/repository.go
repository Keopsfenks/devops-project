package interfaces

import "context"

type IRepository[T any] interface {
	Find(ctx context.Context, filter map[string]interface{}, skip, limit int) ([]*T, error)
	FindOne(ctx context.Context, filter map[string]interface{}) (*T, error)
	InsertOne(ctx context.Context, entity *T) (*T, error)
	ReplaceOne(ctx context.Context, filter map[string]interface{}, entity *T) (*T, error)
	DeleteOne(ctx context.Context, filter map[string]interface{}) (bool, error)
	SoftDeleteOne(ctx context.Context, filter map[string]interface{}, entity *T) (*T, error)
	Exists(ctx context.Context, filter map[string]interface{}) (bool, error)
}

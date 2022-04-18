package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type IClient interface {
	DB() *mongo.Database
	Ping(ctx context.Context) error
	WithTX(ctx context.Context, fn func(context.Context) error) error
	Close() error

	Create(ctx context.Context, rec interface{}) error
	Update(ctx context.Context, rec interface{}) error
	UpdateWhere(ctx context.Context, rec interface{}, filter Filter) (int64, error)
	Upsert(ctx context.Context, rec interface{}) error
	Delete(ctx context.Context, rec interface{}) error
	DeleteWhere(ctx context.Context, rec interface{}, filter Filter) (int64, error)
	GetByID(ctx context.Context, rec interface{}) error
	Find(ctx context.Context, rec interface{}, filter SearchFilter) error
}

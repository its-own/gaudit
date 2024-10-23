package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// NoSql interface wraps the database
type NoSql interface {
	Ping(ctx context.Context) error
	Disconnect(ctx context.Context) error
	EnsureIndices(ctx context.Context, tab string, index []Index) error
	DropIndices(ctx context.Context, tab string, index []Index) error
	Insert(ctx context.Context, tab string, v interface{}) error
	Update(ctx context.Context, col string, filter interface{}, data interface{}) error
	InsertMany(ctx context.Context, tab string, v []interface{}) error
	Count(ctx context.Context, col string, q interface{}) (int64, error)
	List(ctx context.Context, tab string, filter interface{}, skip, limit int64, v interface{}, sort ...interface{}) error
	FindOne(ctx context.Context, tab string, filter interface{}, v interface{}, sort ...interface{}) error
	PartialUpdateMany(ctx context.Context, col string, filter interface{}, data interface{}) error
	PartialUpdateManyByQuery(ctx context.Context, col string, filter interface{}, query UnorderedDbQuery) error
	BulkUpdate(ctx context.Context, col string, models []mongo.WriteModel) error
	Aggregate(ctx context.Context, col string, q []interface{}, v interface{}) error
	AggregateWithDiskUse(ctx context.Context, col string, q []interface{}, v interface{}) error
	Distinct(ctx context.Context, col, field string, q interface{}, v interface{}) error
	DeleteMany(ctx context.Context, col string, filter interface{}) error
}

// Index holds database index
type Index struct {
	Name        string
	Keys        []IndexKey
	Unique      *bool
	Sparse      *bool
	ExpireAfter *time.Duration
}

type IndexKey struct {
	Key string
	Asc interface{}
}
type UnorderedDbQuery bson.M

type BulkWriteModel mongo.WriteModel

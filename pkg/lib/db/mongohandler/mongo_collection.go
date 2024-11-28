package mongohandler

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection interface {
	CreateManyIndexes(ctx context.Context, models []mongo.IndexModel) ([]string, error)
	CountDocuments(ctx context.Context, filter interface{}, opt ...*options.CountOptions) (int64, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

type collection struct {
	collection *mongo.Collection
}

func (c *collection) CreateManyIndexes(ctx context.Context, models []mongo.IndexModel) ([]string, error) {
	return c.collection.Indexes().CreateMany(ctx, models)
}

func (c *collection) CountDocuments(ctx context.Context, filter interface{}, opt ...*options.CountOptions) (int64, error) {
	return c.collection.CountDocuments(ctx, filter, opt...)
}

func (c *collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return c.collection.Find(ctx, filter, opts...)
}

func (c *collection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return c.collection.FindOne(ctx, filter, opts...)
}
func (c *collection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return c.collection.FindOneAndUpdate(ctx, filter, update, opts...)
}

func (c *collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.collection.InsertOne(ctx, document, opts...)
}

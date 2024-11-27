package repository

import (
	"context"

	"github.com/nsaltun/userapi/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, id string, user *model.User) (*model.User, error)
	ListByFilter(ctx context.Context, filter bson.M, limit int, offset int) ([]model.User, int64, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*model.User, error)
}

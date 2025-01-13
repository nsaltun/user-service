package user

import (
	"github.com/nsaltun/userapi/internal/model"
)

const (
	// DefaultLimit is the default limit for pagination
	DefaultLimit = 20
)

type CreateUserRequest struct {
	*model.User
}

type CreateUserResponse struct {
	*model.User
}

type UpdateUserByIdRequest struct {
	*model.User
}

type UpdateUserByIdResponse struct {
	*model.User
}

type ListUsersByFilterRequest struct {
	Limit  int
	Offset int
	*model.UserFilter
}

type ListUsersByFilterResponse struct {
	*model.Pagination
}

type DeleteUserByIdRequest struct {
	Id string `json:"id"`
}

type DeleteUserByIdResponse struct{}

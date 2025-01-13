package user

import (
	"context"
	"net/http"

	"github.com/nsaltun/userapi/internal/service"
)

// UserHandler is an interface for http handler methods for user operations
type UserHandler interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, int, error)
	UpdateUserById(context.Context, *UpdateUserByIdRequest) (*UpdateUserByIdResponse, int, error)
	DeleteUserById(context.Context, *DeleteUserByIdRequest) (*DeleteUserByIdResponse, int, error)
	ListUsers(ctx context.Context, req *ListUsersByFilterRequest) (*ListUsersByFilterResponse, int, error)
}

// Implementor of user handler
type userHandler struct {
	userService service.UserService
}

// NewUserHandler returns an instance of user handler to be able to use it in http router
func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService}
}

// CreateUser is handling user creation. If there is no error it returns
// json data with created record `ID` and `CreatedAt` and `Status` info with 200 http status code.
//
// Validation as required for `firstName`,`email`,`nickName`,`country` fields.
// It returns Http 400 error for validation error.
//
// There is unique constraint for `nickName` and `email` so when new data added with one of those values
// it returns Http 409(conflict) error.
//
// If error occurs it returns structured json data which composed with error code and error message
func (u *userHandler) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, int, error) {
	createdUser, err := u.userService.CreateUser(ctx, req.User)
	if err != nil {
		return nil, 0, err
	}

	return &CreateUserResponse{createdUser}, http.StatusCreated, nil
}

// UpdateUserById is handling user update. If there is no error it returns updated user in json format with 200 http status code
//
// Getting id from path. Getting payload from request body and decoding to user model.
//
// NOTE: It is updating values without comparing if it's changed or not or empty.
// So be careful to send the same data for the fields you don't want to change.
//
// If error occurs it returns structured json data which composed with error code and error message
func (u *userHandler) UpdateUserById(ctx context.Context, req *UpdateUserByIdRequest) (*UpdateUserByIdResponse, int, error) {
	updatedUser, err := u.userService.UpdateUserById(ctx, req.Id, *req.User)
	if err != nil {
		return nil, 0, err
	}
	return &UpdateUserByIdResponse{updatedUser}, http.StatusOK, nil
}

// ListUsers is handling user filtration. If there is no error it returns paginated and filtered user data with 200 http status code.
//
// Getting `limit` and `offset` from the path. Getting filter payload from request body and decoding it into UserFilter model.
// NOTE: If limit and offset not provided it is set to default values which is 20(limit) and (0)offset.
// Also paylaod is not mandatory. If nothing provided in payload it will return all users according to limit and offset.
//
// It only returns active users (user.status=1)
//
// If error occurs it returns structured json data which composed with error code and error message
func (u *userHandler) ListUsers(ctx context.Context, req *ListUsersByFilterRequest) (*ListUsersByFilterResponse, int, error) {
	paginatedData, err := u.userService.ListUsers(ctx, *req.UserFilter, req.Limit, req.Offset)
	if err != nil {
		return nil, 0, err
	}

	return &ListUsersByFilterResponse{paginatedData}, http.StatusOK, nil
}

// DeleteUserById is handling user deletion. If there is no error it returns empty response and HTTP 200 status code
//
// Getting id from path.
//
// It actually updates status field in DB to `2` which means "Inactive".
//
// If error occurs it returns structured json data which composed with error code and error message
func (u *userHandler) DeleteUserById(ctx context.Context, req *DeleteUserByIdRequest) (*DeleteUserByIdResponse, int, error) {
	err := u.userService.DeleteUserById(ctx, req.Id)
	if err != nil {
		return nil, 0, err
	}
	return &DeleteUserByIdResponse{}, http.StatusOK, nil
}

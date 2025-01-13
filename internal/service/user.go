package service

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/nsaltun/userapi/internal/model"
	"github.com/nsaltun/userapi/internal/repository"
	"github.com/nsaltun/userapi/pkg/lib/crypt"
	"github.com/nsaltun/userapi/pkg/lib/errwrap"
	"golang.org/x/crypto/bcrypt"
)

// UserService interface
type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUserById(ctx context.Context, id string, user model.User) (*model.User, error)
	DeleteUserById(ctx context.Context, id string) error
	ListUsers(ctx context.Context, userFilter model.UserFilter, limit int, offset int) (*model.Pagination, error)
}

// userService implementor
type userService struct {
	userRepository repository.UserRepository
}

// NewUserService returns new instance of UserService to use it's methods
func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository}
}

// CreateUser calling relevant repository method to create user.
//
// Error cases:
//
// - Returns internal error when hash is faulty or error from repository other than conflict
//
// - Returns BadRequest when password is too long.
//
// - Returns Conflict error when unique index constraint violated.
//
// Returns created user with ID,CreatedAt,UpdatedAt,Status when operation is successful.
func (u *userService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	hashedPwd, err := crypt.HashPassword(user.Password)
	if err != nil {
		code := http.StatusInternalServerError
		message := "unexpected error"
		if err == bcrypt.ErrPasswordTooLong {
			code = http.StatusBadRequest
			message = "password is too long"
		}
		return nil, errwrap.NewError(message, strconv.Itoa(code)).SetHttpCode(code)
	}

	user.Password = hashedPwd
	err = u.userRepository.Create(ctx, user)
	if err != nil {
		slog.Info("error while creating user.", slog.Any("error", err.Error()))
		return nil, err
	}

	return user, nil
}

// UpdateUserById is calling relevant repository method to update user.
//
// Error cases:
//
// - Returns internal error when hash is faulty.
//
// - Returns BadRequest when password is too long.
//
// - Returns Conflict error when unique index constraint violated.
//
// Returns created user with ID,CreatedAt,UpdatedAt,Status when operation is successful.
func (u *userService) UpdateUserById(ctx context.Context, id string, user model.User) (*model.User, error) {
	user.Id = ""
	updatedUser, err := u.userRepository.Update(ctx, id, &user)
	if err != nil {
		slog.Info("error from repository", slog.Any("error", err.Error()))
		return nil, err
	}

	return updatedUser, nil
}

// DeleteUserByID updates user status to `Inactive(2)`
//
// Returns NotFound error if record not found
func (u *userService) DeleteUserById(ctx context.Context, id string) error {
	return u.userRepository.Delete(ctx, id)
}

// ListUsers lists users with filter and pagination
func (u *userService) ListUsers(ctx context.Context, userFilter model.UserFilter, limit int, offset int) (*model.Pagination, error) {
	users, totalCount, err := u.userRepository.ListByFilter(ctx, userFilter.ToBson(), limit, offset)
	if err != nil {
		slog.Info("error from DB while getting list of users", slog.Any("error", err.Error()))
		return nil, err
	}

	//TODO: validate List user request

	// Determine if there are next and previous pages
	hasNext := int64(offset+limit) < totalCount
	hasPrevious := offset > 0

	// Construct the Pagination response
	pagination := &model.Pagination{
		TotalRecords: totalCount,
		Limit:        limit,
		Offset:       offset,
		HasNext:      hasNext,
		HasPrevious:  hasPrevious,
		Items:        users,
	}

	return pagination, nil
}

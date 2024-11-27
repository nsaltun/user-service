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

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUserById(ctx context.Context, id string, user model.User) (*model.User, error)
	DeleteUserById(ctx context.Context, id string) error
	ListUsers(ctx context.Context, userFilter model.UserFilter, limit int, offset int) (*model.Pagination, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository}
}

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

func (u *userService) UpdateUserById(ctx context.Context, id string, user model.User) (*model.User, error) {
	updatedUser, err := u.userRepository.Update(ctx, id, &user)
	if err != nil {
		slog.Info("error from repository", slog.Any("error", err.Error()))
		return nil, err
	}

	return updatedUser, nil
}

func (u *userService) DeleteUserById(ctx context.Context, id string) error {
	return u.userRepository.Delete(ctx, id)
}

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

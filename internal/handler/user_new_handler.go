package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/nsaltun/userapi/internal/model"
	"github.com/nsaltun/userapi/pkg/lib/errwrap"
)

type CreateUserRequest struct {
	*model.User
}

type CreateUserResponse struct {
	*model.User
}

func (u *userHandler) CreateUserCpy(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, int, error) {
	createdUser, err := u.userService.CreateUser(ctx, req.User)
	if err != nil {
		return nil, 0, err
	}

	return &CreateUserResponse{createdUser}, http.StatusCreated, nil
}

func (req CreateUserRequest) Validate() error {
	validationErrs := []string{}
	if req.FirstName == "" {
		validationErrs = append(validationErrs, "firstName can't be empty")
	}
	if req.Email == "" {
		validationErrs = append(validationErrs, "email can't be empty")
	}
	if req.NickName == "" {
		validationErrs = append(validationErrs, "nickName can't be empty")
	}
	if req.Country == "" {
		validationErrs = append(validationErrs, "country can't be empty")
	}

	if len(validationErrs) > 0 {
		errMsg := strings.Join(validationErrs, ";;")
		return errwrap.ErrBadRequest.SetMessage(errMsg)
	}

	return nil
}

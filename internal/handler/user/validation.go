package user

import (
	"strings"

	"github.com/nsaltun/userapi/pkg/lib/errwrap"
)

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

func (req UpdateUserByIdRequest) Validate() error {
	validationErrs := []string{}

	if req.Id == "" {
		validationErrs = append(validationErrs, "id can't be empty")
	}

	if len(validationErrs) > 0 {
		errMsg := strings.Join(validationErrs, ";;")
		return errwrap.ErrBadRequest.SetMessage(errMsg)
	}
	return nil
}

func (req ListUsersByFilterRequest) Validate() error {
	validationErrs := []string{}
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Limit < 0 {
		validationErrs = append(validationErrs, "limit can't be negative")
	}
	if req.Offset < 0 {
		validationErrs = append(validationErrs, "offset can't be negative")
	}
	if req.UserFilter == nil {
		validationErrs = append(validationErrs, "user filter can't be nil")
	}

	if len(validationErrs) > 0 {
		errMsg := strings.Join(validationErrs, ";;")
		return errwrap.ErrBadRequest.SetMessage(errMsg)
	}
	return nil
}

func (req DeleteUserByIdRequest) Validate() error {
	validationErrs := []string{}

	if req.Id == "" {
		validationErrs = append(validationErrs, "id can't be empty")
	}

	if len(validationErrs) > 0 {
		errMsg := strings.Join(validationErrs, ";;")
		return errwrap.ErrBadRequest.SetMessage(errMsg)
	}
	return nil
}

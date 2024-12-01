package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/nsaltun/userapi/internal/model"
	"github.com/nsaltun/userapi/internal/service"
	"github.com/nsaltun/userapi/pkg/lib/errwrap"
	"github.com/nsaltun/userapi/pkg/lib/middleware"
)

// UserHandler is an interface for http handler methods for user operations
type UserHandler interface {
	CreateUser(c *middleware.HttpContext) error
	UpdateUserById(c *middleware.HttpContext) error
	DeleteUserById(c *middleware.HttpContext) error
	ListUsers(c *middleware.HttpContext) error
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
func (u *userHandler) CreateUser(c *middleware.HttpContext) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		slog.Info("error while unmarshalling request body.", slog.Any("error", err.Error()))
		return errorResp(c, http.StatusBadRequest, errwrap.ErrBadRequest.SetMessage("invalid payload format"))
	}

	//TODO: validate create user request
	validationErrs := []string{}
	if user.FirstName == "" {
		validationErrs = append(validationErrs, "firstName can't be empty")
	}
	if user.Email == "" {
		validationErrs = append(validationErrs, "email can't be empty")
	}
	if user.NickName == "" {
		validationErrs = append(validationErrs, "nickName can't be empty")
	}
	if user.Country == "" {
		validationErrs = append(validationErrs, "country can't be empty")
	}

	if len(validationErrs) > 0 {
		errMsg := strings.Join(validationErrs, ";;")
		return errorResp(c, http.StatusBadRequest, errwrap.ErrBadRequest.SetMessage(errMsg))
	}

	createdUser, err := u.userService.CreateUser(c.Request.Context(), &user)
	if err != nil {
		return errorRespWithMapping(c, err)
	}

	return successResp(c, http.StatusCreated, createdUser)
}

// UpdateUserById is handling user update. If there is no error it returns updated user in json format with 200 http status code
//
// Getting id from path. Getting payload from request body and decoding to user model.
//
// NOTE: It is updating values without comparing if it's changed or not or empty.
// So be careful to send the same data for the fields you don't want to change.
//
// If error occurs it returns structured json data which composed with error code and error message
func (u *userHandler) UpdateUserById(c *middleware.HttpContext) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		slog.Info("error while unmarshalling request body.", slog.Any("error", err.Error()))
		return errorResp(c, http.StatusBadRequest, errwrap.ErrBadRequest.SetMessage(err.Error()))
	}

	//TODO: validate update user request

	id := c.Param("id")
	if len(id) == 0 {
		slog.Info("Id is empty")
		return errorResp(c, http.StatusBadRequest, errwrap.ErrBadRequest.SetMessage("Please provide id in endpoint"))
	}

	updatedUser, err := u.userService.UpdateUserById(c.Request.Context(), id, user)
	if err != nil {
		return errorRespWithMapping(c, err)
	}

	return successResp(c, http.StatusOK, updatedUser)
}

// DeleteUserById is handling user deletion. If there is no error it returns empty response and HTTP 200 status code
//
// Getting id from path.
//
// It actually updates status field in DB to `2` which means "Inactive".
//
// If error occurs it returns structured json data which composed with error code and error message
func (u *userHandler) DeleteUserById(c *middleware.HttpContext) error {
	id := c.Param("id")
	if id == "" {
		return errorResp(c, http.StatusBadRequest, errwrap.ErrBadRequest.SetMessage("id not provided"))
	}

	err := u.userService.DeleteUserById(c.Request.Context(), id)
	if err != nil {
		return errorRespWithMapping(c, err)
	}

	return successResp(c, http.StatusOK, nil)
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
func (u *userHandler) ListUsers(c *middleware.HttpContext) error {
	var userReq model.UserFilter
	if err := c.BodyParser(&userReq); err != nil {
		slog.Info("error while unmarshalling request body.", slog.Any("error", err.Error()))
		return errorResp(c, http.StatusBadRequest, errwrap.ErrBadRequest.SetMessage(err.Error()))
	}

	// Pagination settings (using limit and offset)
	limit := c.QueryInt("limit", 20)  // Default to 20 if not provided
	offset := c.QueryInt("offset", 0) // Default to 0 if not provided

	paginatedData, err := u.userService.ListUsers(c.Request.Context(), userReq, limit, offset)
	if err != nil {
		return errorRespWithMapping(c, err)
	}

	return successResp(c, http.StatusOK, paginatedData)
}

// successResp is taking status and data as input and writing as json data into http response writer.
func successResp(c *middleware.HttpContext, httpStatus int, data interface{}) error {
	if data == nil {
		data = struct{}{}
	}
	return c.JSON(httpStatus, data)
}

// errorResp is taking httpstatus and error as input and writing as json data into http response writer.
func errorResp(c *middleware.HttpContext, httpstatus int, err errwrap.IError) error {
	return c.JSON(httpstatus, err.ErrorResp())
}

// errorResp is taking error as input and mapping it according to error code and writing as json data into http response writer.
func errorRespWithMapping(c *middleware.HttpContext, err error) error {
	iError, ok := err.(errwrap.IError)
	if !ok {
		return c.JSON(http.StatusInternalServerError, errwrap.ErrInternal.ErrorResp())
	}
	return c.JSON(iError.HttpCode(), iError.ErrorResp())
}

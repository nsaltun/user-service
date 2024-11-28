package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/nsaltun/userapi/internal/model"
	"github.com/nsaltun/userapi/internal/service"
	"github.com/nsaltun/userapi/pkg/lib/errwrap"
	"github.com/nsaltun/userapi/pkg/lib/middlewares/httpwrap"
)

type UserHandler interface {
	CreateUser(c *httpwrap.HttpContext) error
	UpdateUserById(c *httpwrap.HttpContext) error
	DeleteUserById(c *httpwrap.HttpContext) error
	ListUsers(c *httpwrap.HttpContext) error
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService}
}

func (u *userHandler) CreateUser(c *httpwrap.HttpContext) error {
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

func (u *userHandler) UpdateUserById(c *httpwrap.HttpContext) error {
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

func (u *userHandler) DeleteUserById(c *httpwrap.HttpContext) error {
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

func (u *userHandler) ListUsers(c *httpwrap.HttpContext) error {
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

func successResp(c *httpwrap.HttpContext, httpStatus int, data interface{}) error {
	if data == nil {
		data = struct{}{}
	}
	return c.JSON(httpStatus, data)
}

func errorResp(c *httpwrap.HttpContext, httpstatus int, err errwrap.IError) error {
	return c.JSON(httpstatus, err.ErrorResp())
}

func errorRespWithMapping(c *httpwrap.HttpContext, err error) error {
	iError, ok := err.(errwrap.IError)
	if !ok {
		return c.JSON(http.StatusInternalServerError, errwrap.ErrInternal.ErrorResp())
	}
	return c.JSON(iError.HttpCode(), iError.ErrorResp())
}

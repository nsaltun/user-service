package handler

import (
	"log/slog"
	"net/http"
	"strconv"

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
	var user *model.User
	if err := c.BodyParser(user); err != nil {
		slog.Info("error while unmarshalling request body.", slog.Any("error", err.Error()))
		return c.JSON(http.StatusBadRequest, errwrap.ErrorResponse{
			Code:    strconv.Itoa(http.StatusBadRequest),
			Message: err.Error(),
		})
	}

	//TODO: validate create user request

	createdUser, err := u.userService.CreateUser(c.Request.Context(), user)
	if err != nil {
		iError, ok := err.(errwrap.IError)
		if !ok {
			return c.JSON(http.StatusInternalServerError, errwrap.ErrInternal.ErrorResp())
		}
		return c.JSON(iError.HttpCode(), iError.ErrorResp())
	}

	return c.JSON(int(http.StatusCreated), createdUser)
}

func (u *userHandler) UpdateUserById(c *httpwrap.HttpContext) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		slog.Info("error while unmarshalling request body.", slog.Any("error", err.Error()))
		return c.JSON(http.StatusBadRequest, errwrap.ErrBadRequest.SetMessage(err.Error()).ErrorResp())
	}

	//TODO: validate update user request

	id := c.Param("id")
	if len(id) == 0 {
		slog.Info("Id is empty")
		return c.JSON(http.StatusBadRequest, errwrap.ErrBadRequest.SetMessage("Please provide id in endpoint").ErrorResp())
	}

	updatedUser, err := u.userService.UpdateUserById(c.Request.Context(), id, user)
	if err != nil {
		iError, ok := err.(errwrap.IError)
		if !ok {
			return c.JSON(http.StatusInternalServerError, errwrap.ErrInternal.ErrorResp())
		}
		return c.JSON(iError.HttpCode(), iError.ErrorResp())
	}

	return c.JSON(int(http.StatusOK), updatedUser)
}

func (u *userHandler) DeleteUserById(c *httpwrap.HttpContext) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, errwrap.ErrBadRequest.SetMessage("id not provided").ErrorResp())
	}

	err := u.userService.DeleteUserById(c.Request.Context(), id)
	if err != nil {
		iError, ok := err.(errwrap.IError)
		if !ok {
			return c.JSON(http.StatusInternalServerError, errwrap.ErrInternal.ErrorResp())
		}
		return c.JSON(iError.HttpCode(), iError.ErrorResp())
	}

	return c.JSON(http.StatusOK, struct{}{})
}

func (u *userHandler) ListUsers(c *httpwrap.HttpContext) error {
	var userReq model.UserFilter
	if err := c.BodyParser(&userReq); err != nil {
		slog.Info("error while unmarshalling request body.", slog.Any("error", err.Error()))
		return c.JSON(http.StatusBadRequest, errwrap.ErrBadRequest.SetMessage(err.Error()).ErrorResp())
	}

	// Pagination settings (using limit and offset)
	limit := c.QueryInt("limit", 20)  // Default to 20 if not provided
	offset := c.QueryInt("offset", 0) // Default to 0 if not provided

	paginatedData, err := u.userService.ListUsers(c.Request.Context(), userReq, limit, offset)
	if err != nil {
		iError, ok := err.(errwrap.IError)
		if !ok {
			return c.JSON(http.StatusInternalServerError, errwrap.ErrInternal.ErrorResp())
		}
		return c.JSON(iError.HttpCode(), iError.ErrorResp())
	}

	return c.JSON(http.StatusOK, paginatedData)
}

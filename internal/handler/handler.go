package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nsaltun/userapi/pkg/lib/errwrap"
	"github.com/nsaltun/userapi/pkg/lib/middleware"
)

type Request interface {
	Validate() error
}

type Response any

// HandlerFunc is a function type that takes a context and a request and returns a response, a success status code and an error.
//
// success status code is processing only for the success response. If the response is successful, the status code will be set to the response.
type HandlerFunc[Request any, Response any] func(context.Context, *Request) (*Response, int, error)

func Serve[I Request, O Response](h HandlerFunc[I, O]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req I
		if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
			return err
		}

		if err := c.ParamsParser(&req); err != nil {
			return err
		}

		if err := c.QueryParser(&req); err != nil {
			return err
		}

		if err := req.Validate(); err != nil {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}

		ctx := c.UserContext()
		resp, successStatusCode, err := h(ctx, &req)
		if err != nil {

			return errorRespWithMapping(err)
		}

		c.Status(successStatusCode)
		return c.JSON(resp)
	}
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
func errorRespWithMapping(err error) *fiber.Error {
	iError, ok := err.(errwrap.IError)
	if !ok {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return fiber.NewError(iError.HttpCode(), iError.Error())
}

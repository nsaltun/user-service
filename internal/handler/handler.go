package handler

import (
	"context"
	"net/http"

	"github.com/nsaltun/userapi/pkg/lib/errwrap"
	"github.com/nsaltun/userapi/pkg/lib/middleware"
)

type Validation interface {
	Validate() error
}

type Request interface {
	Validation
}
type Response any

// HandlerFunc is a function type that takes a context and a request and returns a response, a success status code and an error.
//
// success status code is processing only for the success response. If the response is successful, the status code will be set to the response.
type HandlerFunc[Request any, Response any] func(context.Context, *Request) (*Response, int, error)

func Serve[I Request, O Response](h HandlerFunc[I, O]) middleware.CustomHandler {
	return func(c *middleware.HttpContext) error {
		req := new(I)
		if err := c.BodyParser(req); err != nil {
			return errorResp(c, http.StatusBadRequest, errwrap.ErrBadRequest.SetMessage("invalid payload format"))
		}

		if err := (*req).Validate(); err != nil {
			return errorResp(c, http.StatusBadRequest, errwrap.NewFromError(err))
		}

		ctx := c.Request.Context()
		resp, successStatusCode, err := h(ctx, req)
		if err != nil {
			return errorRespWithMapping(c, err)
		}
		return successResp(c, successStatusCode, resp)
	}
}

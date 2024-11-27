package errwrap

import "net/http"

var (
	ErrBadRequest = NewError("invalid argument", "400").SetHttpCode(http.StatusBadRequest)
	ErrNotFound   = NewError("resource not found", "404").SetHttpCode(http.StatusNotFound)
	ErrConflict   = NewError("already exists", "409").SetHttpCode(http.StatusConflict)
	ErrInternal   = NewError("internal server error", "500").SetHttpCode(http.StatusInternalServerError)
)

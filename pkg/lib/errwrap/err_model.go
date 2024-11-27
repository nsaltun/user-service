package errwrap

import "fmt"

type IError interface {
	error
	SetMessage(msg string) IError
	SetHttpCode(code int) IError
	HttpCode() int
	ErrorResp() ErrorResponse
}

type errorWrapper struct {
	message  string
	code     string
	httpCode int
}

type ErrorResponse struct {
	Message string
	Code    string
}

func NewError(msg string, code string) IError {
	return &errorWrapper{
		message: msg,
		code:    code,
	}
}

func (e *errorWrapper) SetMessage(msg string) IError {
	newErr := e.clone()
	newErr.message = msg
	return newErr
}

func (e *errorWrapper) SetHttpCode(code int) IError {
	newErr := e.clone()
	newErr.httpCode = code
	return newErr
}

func (e *errorWrapper) HttpCode() int {
	return e.httpCode
}

func (e *errorWrapper) ErrorResp() ErrorResponse {
	return ErrorResponse{
		Code:    e.code,
		Message: e.message,
	}
}

func (e *errorWrapper) Error() string {
	return fmt.Sprintf("%s code:%s", e.message, e.code)
}

func (e *errorWrapper) clone() *errorWrapper {
	if e == nil {
		return nil
	}
	return &errorWrapper{
		code:     e.code,
		httpCode: e.httpCode,
		message:  e.message,
	}
}

package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

// HttpContext is a custom wrapper around http.ResponseWriter and http.Request
type HttpContext struct {
	Response http.ResponseWriter
	Request  *http.Request
	UserID   string // custom field (e.g., for auth)
}

// Middleware type that accepts and returns a function with HttpContext
type Middleware func(next CustomHandler) CustomHandler

// CustomHandler is the handler signature using HttpContext
type CustomHandler func(*HttpContext) error

// MiddlewareRunner chains middlewares and passes HttpContext to the final handler
func MiddlewareRunner(handler CustomHandler, middlewares ...Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize HttpContext
		ctx := &HttpContext{Response: w, Request: r}

		// Apply middlewares recursively
		final := handler
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}

		// Execute the final handler with HttpContext
		final(ctx)
	}
}

// Param retrieves a parameter from the URL by name.
func (c *HttpContext) Param(key string) string {
	return c.Request.PathValue(key)
}

// BodyParser parses the JSON body into a provided struct
func (c *HttpContext) BodyParser(model interface{}) error {
	if c.Request.Body == nil {
		return errors.New("request body is empty")
	}
	defer c.Request.Body.Close()
	return json.NewDecoder(c.Request.Body).Decode(model)
}

// JSON sends a JSON response with the given status code
func (c *HttpContext) JSON(status int, payload interface{}) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.Header().Set("statusCode", fmt.Sprint(status))
	c.Response.WriteHeader(status)
	json.NewEncoder(c.Response).Encode(payload)
	return nil
}

func (c *HttpContext) QueryInt(key string, defaultValue int) int {
	val := c.Request.URL.Query().Get(key)
	if val == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		slog.Info("value in query is not an integer.", slog.Any("queryValue", val))
		return defaultValue
	}

	return intVal
}

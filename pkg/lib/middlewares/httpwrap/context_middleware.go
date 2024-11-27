package httpwrap

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

type HandlerFunc func(*HttpContext) error

// HttpContext is a custom wrapper around http.ResponseWriter and http.Request
type HttpContext struct {
	Response http.ResponseWriter
	Request  *http.Request
}

// NewContext creates a new Context instance
func NewContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	return &HttpContext{
		Response: w,
		Request:  r,
	}
}

// ContextMiddleware injects the custom Context and wraps the custom handler
func ContextMiddleware(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create the custom context
		ctx := NewContext(w, r)

		// Call the custom handler with the context
		next(ctx)
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

package fiber_middleware

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Response structure for consistent API responses
type APIResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponseMiddleware is a Fiber middleware for mapping and logging responses
func ResponseMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now() // Record the start time

		// Proceed with the next handler
		err := c.Next()

		// Calculate response time
		duration := time.Since(start)
		statusCode := c.Response().StatusCode()
		level := slog.LevelInfo

		// Initialize a response wrapper
		response := APIResponse{
			Status:  "success",
			Code:    statusCode,
			Message: "",
			Data:    nil,
		}

		// Handle errors
		if err != nil {
			response.Status = "error"
			if fiberErr, ok := err.(*fiber.Error); ok {
				// Fiber error
				response.Code = fiberErr.Code
				response.Message = fiberErr.Message
				level = slog.LevelWarn
			} else {
				// Generic error
				response.Code = fiber.StatusInternalServerError
				response.Message = err.Error()
				level = slog.LevelError
			}
		} else {
			// Parse the original response body
			if len(c.Response().Body()) > 0 {
				if err := json.Unmarshal(c.Response().Body(), &response.Data); err != nil {
					slog.Warn("Failed to parse response body", "error", err)
					response.Data = string(c.Response().Body())
				}
			} else {
				response.Data = fiber.Map{"message": "Request completed successfully"}
			}
		}

		// Log the request and response
		slog.Log(c.Context(), level, "request completed",
			"method", c.Method(),
			"path", c.Path(),
			"status", response.Status,
			"statusCode", response.Code,
			"duration", duration,
		)

		// Send the final wrapped response
		c.Response().Reset()

		// Set the content type and send the JSON response
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(statusCode).JSON(response)
	}
}

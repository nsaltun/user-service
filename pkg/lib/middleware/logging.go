package middleware

import (
	"log/slog"
	"time"
)

// LoggingMiddleware logs request details
func LoggingMiddleware(next CustomHandler) CustomHandler {
	return func(ctx *HttpContext) error {
		start := time.Now()

		// Log the incoming request
		slog.Info("Incoming Request",
			"method", ctx.Request.Method,
			"url", ctx.Request.URL.String(),
			// "headers", ctx.Request.Header,
		)
		next(ctx)
		slog.Info("Outgoing Response",
			"status", ctx.Response.Header().Get("statusCode"),
			"duration(ns)", time.Since(start),
		)
		return nil
	}
}

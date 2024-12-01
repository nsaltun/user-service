package router

import (
	"net/http"

	"github.com/nsaltun/userapi/internal/handler"
	"github.com/nsaltun/userapi/pkg/lib/health"
	"github.com/nsaltun/userapi/pkg/lib/middleware"
)

// NewRouter adds endpoint patterns with handlers and middlewares into http.ServeMux and returns that mux
func NewRouter(userHandler handler.UserHandler, health health.HealthCheck) http.Handler {
	mux := http.NewServeMux()
	// Register routes using the custom context handler
	mux.HandleFunc("POST /users", middleware.MiddlewareRunner(userHandler.CreateUser, middleware.LoggingMiddleware, middleware.AuthMiddleware))
	mux.HandleFunc("PUT /users/{id}", middleware.MiddlewareRunner(userHandler.UpdateUserById, middleware.LoggingMiddleware, middleware.AuthMiddleware))
	mux.HandleFunc("DELETE /users/{id}", middleware.MiddlewareRunner(userHandler.DeleteUserById, middleware.LoggingMiddleware, middleware.AuthMiddleware))
	mux.HandleFunc("POST /users/filter", middleware.MiddlewareRunner(userHandler.ListUsers, middleware.LoggingMiddleware, middleware.AuthMiddleware))

	//healthchecker
	mux.HandleFunc("GET /health", health.HealthCheckHandler())

	return mux
}

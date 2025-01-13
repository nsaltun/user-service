package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nsaltun/userapi/internal/handler"
	"github.com/nsaltun/userapi/internal/handler/user"
	"github.com/nsaltun/userapi/pkg/lib/health"
	"github.com/nsaltun/userapi/pkg/lib/middleware/fiber_middleware"
)

func NewFiberRouter(app *fiber.App, userHandler user.UserHandler, health health.HealthCheck) {
	// Use the response middleware
	userApi := app.Group("/api/users")
	userApi.Use(fiber_middleware.ResponseMiddleware())
	userApi.Post("", handler.Serve(userHandler.CreateUser))
	userApi.Put("/:id", handler.Serve(userHandler.UpdateUserById))
	userApi.Post("/filter", handler.Serve(userHandler.ListUsers))
	userApi.Delete("/:id", handler.Serve(userHandler.DeleteUserById))
}

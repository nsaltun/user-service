package main

import (
	"log"
	"log/slog"

	"github.com/nsaltun/userapi/internal/handler/user"
	"github.com/nsaltun/userapi/internal/repository"
	"github.com/nsaltun/userapi/internal/router"
	"github.com/nsaltun/userapi/internal/service"
	"github.com/nsaltun/userapi/pkg/lib/db/mongohandler"
	"github.com/nsaltun/userapi/pkg/lib/health"
	"github.com/nsaltun/userapi/pkg/lib/httpserver"
	"github.com/nsaltun/userapi/pkg/lib/logging"
)

func main() {
	logging.InitSlog()
	slog.Info("----USER API----")

	mongodb := mongohandler.New()
	mongodb.InitMongoDB()
	defer mongodb.Disconnect()

	userRepo, err := repository.NewUserRepository(mongodb)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	userSvc := service.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userSvc)

	healthChecker := health.NewHealthCheck(mongodb.HealthChecker())
	// httpHandler := router.NewRouter(userHandler, healthChecker)

	fiberApp := httpserver.NewFiberServer()
	router.NewFiberRouter(fiberApp.App, userHandler, healthChecker)
	fiberApp.Listen()
}

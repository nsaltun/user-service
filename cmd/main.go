package main

import (
	"log"
	"log/slog"

	"github.com/nsaltun/userapi/internal/handler"
	"github.com/nsaltun/userapi/internal/repository"
	"github.com/nsaltun/userapi/internal/router"
	"github.com/nsaltun/userapi/internal/service"
	"github.com/nsaltun/userapi/pkg/lib/db/mongohandler"
	"github.com/nsaltun/userapi/pkg/lib/healthcheck"
	"github.com/nsaltun/userapi/pkg/lib/httpserver"
	"github.com/nsaltun/userapi/pkg/lib/logging"
)

func main() {
	logging.InitSlog()
	slog.Info("----USER API----")

	mongodb := mongohandler.InitMongoDB()
	defer mongodb.Disconnect()

	userRepo, err := repository.NewUserRepository(mongodb)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	mux := router.NewRouter(userHandler)
	health := healthcheck.NewHealthCheck(mongodb)
	health.Init(mux)

	server := httpserver.NewServer(mux)
	server.InitServer()
}

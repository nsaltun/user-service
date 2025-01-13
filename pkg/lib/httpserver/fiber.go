package httpserver

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type FiberServer struct {
	App     *fiber.App
	Address string
}

func NewFiberServer() *FiberServer {
	vi := viper.New()
	vi.AutomaticEnv()
	vi.SetDefault("HOST_ADDRESS", "0.0.0.0")
	vi.SetDefault("PORT", "8080")
	vi.SetDefault("READ_TIMEOUT_IN_SECONDS", 10)
	vi.SetDefault("WRITE_TIMEOUT_IN_SECONDS", 10)
	vi.SetDefault("IDLE_TIMEOUT_IN_SECONDS", 10)
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(vi.GetInt("READ_TIMEOUT_IN_SECONDS")) * time.Second,
		WriteTimeout: time.Duration(vi.GetInt("WRITE_TIMEOUT_IN_SECONDS")) * time.Second,
		IdleTimeout:  time.Duration(vi.GetInt("IDLE_TIMEOUT_IN_SECONDS")) * time.Second,
	})
	return &FiberServer{
		App:     app,
		Address: fmt.Sprintf("%s:%s", vi.Get("HOST_ADDRESS"), vi.Get("PORT")),
	}
}

func (s *FiberServer) Listen() {
	// Channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Run the server in a goroutine
	go func() {
		if err := s.App.Listen(s.Address); err != nil {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	slog.Info(fmt.Sprintf("Server is running on %s", s.Address))

	// Block until a termination signal is received
	<-quit
	slog.Info("Server is shutting down...")

	// Create a timeout context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := s.App.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	slog.Info("Server stopped gracefully")
}

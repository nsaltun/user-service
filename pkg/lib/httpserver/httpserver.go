package httpserver

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

type HttpServer struct {
	port   string
	server *http.Server
}

func NewServer(mux *http.ServeMux) *HttpServer {
	vi := viper.New()
	vi.AutomaticEnv()
	vi.SetDefault("HOST_ADDRESS", "localhost")
	vi.SetDefault("PORT", "3000")
	vi.SetDefault("READ_TIMEOUT_IN_SECONDS", 10)
	vi.SetDefault("WRITE_TIMEOUT_IN_SECONDS", 10)
	vi.SetDefault("IDLE_TIMEOUT_IN_SECONDS", 10)

	return &HttpServer{
		port: vi.GetString("PORT"),
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", vi.Get("HOST_ADDRESS"), vi.Get("PORT")),
			Handler:      mux,
			ReadTimeout:  time.Duration(vi.GetInt("READ_TIMEOUT_IN_SECONDS")) * time.Second,
			WriteTimeout: time.Duration(vi.GetInt("WRITE_TIMEOUT_IN_SECONDS")) * time.Second,
			IdleTimeout:  time.Duration(vi.GetInt("IDLE_TIMEOUT_IN_SECONDS")) * time.Second,
		},
	}

}

func (s *HttpServer) InitServer() {

	// Graceful Shutdown Handling
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
	slog.Info(fmt.Sprintf("Server is running on :%s", s.port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	slog.Info("Server stopped gracefully")
}

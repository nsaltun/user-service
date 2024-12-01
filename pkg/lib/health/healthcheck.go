package health

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/nsaltun/userapi/pkg/lib/db/mongohandler"
	"golang.org/x/sync/errgroup"
)

type HealthCheck interface {
	HealthCheckHandler() http.HandlerFunc
}

type HealthCheckResponse struct {
	Status  string            `json:"status"`
	Details map[string]string `json:"details,omitempty"`
}

type healthCheck struct {
	mongoHealthCheck mongohandler.HealthFn
}

func NewHealthCheck(mongoHealthCheck mongohandler.HealthFn) HealthCheck {
	return &healthCheck{
		mongoHealthCheck,
	}
}

// HealthCheckHandler handles the /health endpoint
func (h *healthCheck) HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		health := HealthCheckResponse{Status: "ok", Details: map[string]string{}}

		// Perform health checks in parallel
		g, ctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			if err := h.mongoHealthCheck(ctx); err != nil {
				health.Status = "unhealthy"
				health.Details["MongoDB"] = "DOWN"
			} else {
				health.Details["MongoDB"] = "UP"
			}
			return nil
		})

		if err := g.Wait(); err != nil {
			health.Status = "unhealthy"
		}

		w.Header().Set("Content-Type", "application/json")
		if health.Status == "unhealthy" {
			slog.ErrorContext(ctx, "service is unhealthy", slog.Any("health", health))
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			slog.DebugContext(ctx, "service is healthy", slog.Any("health", health))
			w.WriteHeader(http.StatusOK)
		}
		json.NewEncoder(w).Encode(health)
	}
}

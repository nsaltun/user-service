package logging

import (
	"log/slog"
	"os"
)

func InitSlog() {
	// Set up a JSON handler for logging
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(jsonHandler))
}

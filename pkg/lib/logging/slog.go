package logging

import (
	"log/slog"
	"os"
	"time"
)

func InitSlog() {
	// Set up a JSON handler for logging
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Replace the default time attribute with a UTC time
			if a.Key == slog.TimeKey {
				a.Value = slog.TimeValue(time.Now().UTC())
			}

			return a
		},
	})

	slog.SetDefault(slog.New(jsonHandler))
}

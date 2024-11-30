package logging

import (
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

var levelMap = map[string]slog.Leveler{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError}

func InitSlog() {
	vi := viper.New()
	vi.AutomaticEnv()
	level := "INFO"
	vi.SetDefault("LOG_LEVEL", level)
	level = vi.GetString("LOG_LEVEL")

	// Set up a JSON handler for logging
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: levelMap[level],
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

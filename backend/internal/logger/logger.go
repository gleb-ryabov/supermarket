package logger

import (
	"errors"
	"log/slog"
	"os"

	"supermarket/internal/config"
)

var (
	errUnsupportedLogLevel = errors.New("unsupported log level")
)

// Init - configure logger.
func Init(logLevel string) (*slog.Logger, error) {
	var slogLevel slog.Level

	switch logLevel {
	case config.LogLevelDebug:
		slogLevel = slog.LevelDebug
	case config.LogLevelInfo:
		slogLevel = slog.LevelInfo
	case config.LogLevelWarn:
		slogLevel = slog.LevelWarn
	case config.LogLevelError:
		slogLevel = slog.LevelError
	default:
		return nil, errUnsupportedLogLevel
	}

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slogLevel,
		}),
	)

	return logger, nil
}

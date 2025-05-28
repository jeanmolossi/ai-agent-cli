package log

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

var Logger *slog.Logger

func Init(level string) error {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return fmt.Errorf("invalid log level %q: %w", level, err)
	}

	format := viper.GetString("log.format")
	if format == "" {
		format = "text"
	}

	var handler slog.Handler

	switch format {
	case "text":
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl})
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: lvl})
	default:
		return fmt.Errorf("format is unsupported, use text or json")
	}

	Logger = slog.New(handler)

	slog.SetDefault(Logger)
	return nil
}

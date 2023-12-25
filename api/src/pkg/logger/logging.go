package logger

import (
	"log/slog"
	"os"
)

func init() {
	slog.Info("initializing logger")

	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	slog.SetDefault(slog.New(h))
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	// if number of args is uneven, assume the last arg is err
	if len(args)%2 != 0 {
		err := args[len(args)-1]
		args = append(args[:len(args)-1], "error", err)
	}

	slog.Error(msg, args...)
}

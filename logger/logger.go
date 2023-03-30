package logger

import "github.com/gookit/slog"

var logger = slog.New()

func Info(args ...any) {
	logger.Info(args)
}

func Error(args ...any) {
	logger.Error(args)
}

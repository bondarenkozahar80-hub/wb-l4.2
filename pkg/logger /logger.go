package logger

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func NewSlogLogger() *slog.Logger {
	consoleHandler := handler.NewConsoleHandler([]slog.Level{
		slog.InfoLevel,
		slog.WarnLevel,
		slog.ErrorLevel,
		slog.FatalLevel,
		slog.PanicLevel,
	})

	lg := slog.NewWithHandlers(consoleHandler)

	slog.PushHandler(consoleHandler)

	return lg
}

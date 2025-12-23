package logger

import (
	"io"
	"log/slog"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
}

type logger struct {
	*slog.Logger
}

func NewLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func (l *logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

func (l *logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

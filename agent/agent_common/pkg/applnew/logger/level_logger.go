package logger

import (
	"io"
	"log/slog"
)

type LogLevel string

const (
	DebugLogLevel LogLevel = "debug"
	InfoLogLevel LogLevel = "info"
	ErrorLogLevel LogLevel = "error"
	WarnLogLevel LogLevel = "warn"
)

type LevelDebugLogger interface {
	Debug(msg string, args ...any)
}
type LevelInfoLogger interface {
	Info(msg string, args ...any)
}
type LevelErrorLogger interface {
	Error(msg string, args ...any)
}
type LevelWarnLogger interface {
	Warn(msg string, args ...any)
}

type LevelLogger interface {
	LevelDebugLogger
	LevelInfoLogger
	LevelErrorLogger
	LevelWarnLogger
}

func convertSlogLevel(level LogLevel) slog.Level {
	switch level {
	case DebugLogLevel:
		return slog.LevelDebug
	case InfoLogLevel:
		return slog.LevelInfo
	default:
		return slog.LevelError
	}
}

func NewSlogLogger(writer io.Writer, level LogLevel) (LevelLogger, error) {
	handle := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		AddSource: true,
		Level: convertSlogLevel(level),
	})

	return slog.New(handle), nil
}
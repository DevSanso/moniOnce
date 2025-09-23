package logger

import "io"

type ILogger interface {
	Log(key string, args ...any) error
}

type ILoggerCloser interface {
	io.Closer
	ILogger
}
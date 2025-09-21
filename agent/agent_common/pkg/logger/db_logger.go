package logger

import "io"

type DbLogger interface {
	Exec(query string, args [][]any) error
}

type DbLoggerCloser interface {
	DbLogger
	io.Closer
}
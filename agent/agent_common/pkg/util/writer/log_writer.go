package writer

import (
	"io"
	"os"
)

type LogWriter struct {
	logFd *os.File
}

func NewLogWriter(path string) (io.WriteCloser, error) {
	var fd *os.File = nil

	if path == "" {
		fd = os.Stdout
	} else {
		f, fErr := os.OpenFile(path, os.O_APPEND | os.O_CREATE, os.FileMode(640))
		if fErr != nil {
			return nil, fErr
		}

		fd = f
	}

	return &LogWriter{
		logFd: fd,
	}, nil
}

func (lw *LogWriter)Write(b []byte) (n int, err error) {
	return lw.logFd.Write(b)
}

func (lw *LogWriter)Close() error {
	return lw.logFd.Close()
}
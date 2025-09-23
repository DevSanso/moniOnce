package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type SizeLimitedWriter struct {
	dir        string
	filename   string
	maxSize    int64 
	current    *os.File
	currentSize int64
	mu         sync.Mutex
}

func NewSizeLimitedWriter(dir, filename string, maxSizeMB int) (*SizeLimitedWriter, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	fullPath := filepath.Join(dir, filename)
	f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &SizeLimitedWriter{
		dir:         dir,
		filename:    filename,
		maxSize:     int64(maxSizeMB) * 1024 * 1024,
		current:     f,
		currentSize: info.Size(),
	}, nil
}

func (w *SizeLimitedWriter) Write(p []byte) (int, error) {
	w.mu.Lock()

	if w.currentSize+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			w.mu.Unlock()
			return 0, err
		}
	}

	n, err := w.current.Write(p)
	w.currentSize += int64(n)
	w.mu.Unlock()
	return n, err
}

func (w *SizeLimitedWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.current != nil {
		err := w.current.Close()
		w.current = nil
		return err
	}
	return nil
}

func (w *SizeLimitedWriter) rotate() error {
	w.current.Close()

	timestamp := time.Now().Format("20060102_150405")
	newName := fmt.Sprintf("%s.%s", w.filename, timestamp)
	err := os.Rename(
		filepath.Join(w.dir, w.filename),
		filepath.Join(w.dir, newName),
	)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(w.dir, w.filename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	w.current = f
	w.currentSize = 0
	return nil
}
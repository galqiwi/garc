package io

import (
	"io"
	"sync"
)

type CountingWriter struct {
	mutex   sync.Mutex
	counter int64
}

func (w *CountingWriter) Write(data []byte) (int, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.counter += int64(len(data))
	return len(data), nil
}

func (w *CountingWriter) Close() error {
	return nil
}

func (w *CountingWriter) BytesWritten() int64 {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.counter
}

var _ io.WriteCloser = &CountingWriter{}

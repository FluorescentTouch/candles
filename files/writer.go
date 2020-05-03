package files

import (
	"errors"
	"os"
	"sync"
)

var errInvalidBytesWrite = errors.New("invalid bytes count written to file")

// Writer represents writer to file.
type Writer struct {
	*sync.Mutex
	file *os.File
}

// NewWriter creates new Writer.
func NewWriter(path string) (*Writer, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &Writer{Mutex: &sync.Mutex{}, file: f}, nil
}

// Close closes the file inside a Writer.
func (w *Writer) Close() {
	w.Lock()
	defer w.Unlock()

	_ = w.file.Close()
}

// WriteString writes a string to file.
func (w *Writer) WriteString(data string) error {
	w.Lock()
	defer w.Unlock()

	n, err := w.file.WriteString(data)
	if err != nil {
		return err
	}

	if n != len([]byte(data)) {
		return errInvalidBytesWrite
	}

	return nil
}

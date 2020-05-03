package pipelines

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/candles/files"
)

// WriterBuilder returns factory of writers.
type WriterBuilder struct{}

// New creates new Writer by given filepath.
func (wb WriterBuilder) New(filepath string) (FileWriter, error) {
	return files.NewWriter(filepath)
}

// Writer describes worker which writes data to corresponding file.
type Writer struct {
	fw   FileWriter
	data <-chan string

	l *logrus.Logger
}

// NewWriter creates new Writer.
func NewWriter(
	fw FileWriter,
	data <-chan string,
	l *logrus.Logger,
) *Writer {
	return &Writer{
		fw:   fw,
		data: data,
		l:    l,
	}
}

func (w *Writer) startWriting(wg *sync.WaitGroup) {
	defer w.fw.Close()

	for data := range w.data {
		err := w.fw.WriteString(data + "\n")
		if err != nil {
			w.l.Errorf("error writing to file: %v", err)
			continue
		}
	}

	wg.Done()
}

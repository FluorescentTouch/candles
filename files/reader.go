package files

import (
	"bufio"
	"os"

	"github.com/sirupsen/logrus"
)

// Reader represents file reader.
type Reader struct {
	fileName string
	file     *os.File
	fileData chan string
	start    chan struct{}

	l *logrus.Logger
}

// NewReader creates new file reader.
func NewReader(filename string, logger *logrus.Logger) (Reader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Reader{}, err
	}

	return Reader{
		fileName: filename,
		file:     f,
		fileData: make(chan string),
		start:    make(chan struct{}),
		l:        logger,
	}, nil
}

// C returns chan which data would be written to.
func (r Reader) C() chan string {
	return r.fileData
}

// Start starts writing data to output chan.
func (r Reader) Init() {
	<-r.start

	s := bufio.NewScanner(r.file)

	for s.Scan() {
		r.fileData <- s.Text()
	}

	if err := s.Err(); err != nil {
		r.l.Errorf("scan error: %v", err)
	}

	close(r.fileData)
}

// StartChan returns chan to receive a start signal.
func (r Reader) StartChan() chan struct{} {
	return r.start
}

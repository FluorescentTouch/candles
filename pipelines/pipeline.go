package pipelines

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/candles/pipelines/candles"
)

var errIntervalAlreadyExists = errors.New("pipeline with provided interval already exists")

var (
	endHours   = 3
	startHours = 10
)

type fileReader interface {
	C() chan string
	StartChan() chan struct{}
	Init()
}

type WritersBuilder interface {
	New(filepath string) (FileWriter, error)
}

type FileWriter interface {
	WriteString(string) error
	Close()
}

// Pipelines describes pipelines aggregator.
type Pipelines struct {
	Start    chan struct{}
	FileDone chan struct{}
	Done     chan struct{}

	r  fileReader
	wb WritersBuilder

	workers []*Worker
	writers []*Writer

	l *logrus.Logger
}

// New creates new pipelines aggregator.
func New(fr fileReader, wb WritersBuilder, l *logrus.Logger) *Pipelines {
	ps := &Pipelines{
		Start:    make(chan struct{}),
		Done:     make(chan struct{}),
		FileDone: make(chan struct{}),

		r:       fr,
		wb:      wb,
		workers: make([]*Worker, 0, 3),
		writers: make([]*Writer, 0, 3),
		l:       l,
	}
	ps.l.Info("Pipelines created")

	return ps
}

// Add adds new pipeline with provided time interval to aggregator.
func (ps *Pipelines) Add(interval int) error {
	for i := range ps.workers {
		if ps.workers[i].interval == interval {
			return errIntervalAlreadyExists
		}
	}

	worker := NewWorker(interval)
	fileName := fmt.Sprintf("candle_%dmin", interval)

	fw, err := ps.wb.New(fileName)
	if err != nil {
		return err
	}

	ps.workers = append(ps.workers, worker)
	ps.writers = append(ps.writers, NewWriter(fw, worker.out, ps.l))

	ps.l.Infof("Pipeline with interval %d added", interval)

	return nil
}

// Init inits pipeline and waits signal for start reading from fileReader.
func (ps *Pipelines) Init() {
	// pipeline stage 3
	go ps.startFileWriters()

	// pipeline stage 2
	for _, w := range ps.workers {
		go w.start()
	}

	go ps.startDataProcess()

	// pipeline stage 1
	go ps.r.Init()

	// wait for start signal
	go func() {
		<-ps.Start
		ps.r.StartChan() <- struct{}{}
	}()
}

// startDataProcess represents start of stage two of pipeline:
// parse trade and sent to workers.
func (ps *Pipelines) startDataProcess() {
	for s := range ps.r.C() {
		tr, err := candles.TradeFromString(s)
		if err != nil {
			ps.l.Errorf("error parsing trade: %s, %v", s, err)
			continue
		}

		if !inWorkingRange(tr.Timestamp) {
			ps.l.Debug("trade is not inside working hours range, skipping", tr)
			continue
		}

		for i := range ps.workers {
			ps.workers[i].in <- tr
		}
	}

	for i := range ps.workers {
		close(ps.workers[i].in)
	}

	ps.FileDone <- struct{}{}
	close(ps.FileDone)
}

// startFileWriters represents start of stage three of pipeline:
// write data to corresponding files.
func (ps *Pipelines) startFileWriters() {
	wg := &sync.WaitGroup{}

	for _, w := range ps.writers {
		wg.Add(1)

		go w.startWriting(wg)
	}

	wg.Wait()

	ps.Done <- struct{}{}
	close(ps.Done)
}

func inWorkingRange(t time.Time) bool {
	return !(t.Hour() >= endHours && t.Hour() < startHours)
}

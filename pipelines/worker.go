package pipelines

import (
	"strings"
	"time"

	"github.com/candles/pipelines/candles"
)

// Worker describes single pipeline with provided time interval.
type Worker struct {
	interval int
	in       chan candles.Trade
	out      chan string

	intervalD     time.Duration
	intervalStart time.Time
	intervalEnd   time.Time
}

// NewWorker creates new pipeline worker with provided interval.
func NewWorker(interval int) *Worker {
	return &Worker{
		interval:  interval,
		intervalD: time.Minute * time.Duration(interval),
		in:        make(chan candles.Trade),
		out:       make(chan string),
	}
}

// incrementInterval checks if trade was in current time interval,
// and if not - increments workers interval-related values.
// Function also handles edge conditions.
func (w *Worker) incrementInterval(trTime time.Time) {
	const hoursInDay = 24

	for trTime.After(w.intervalEnd) || trTime.Equal(w.intervalEnd) {
		newStart := w.intervalStart.Add(w.intervalD)
		if (w.intervalStart == time.Time{} || !inWorkingRange(newStart)) {
			// set start interval for 10:00 of current day if not already set.
			w.intervalStart = trTime.Truncate(time.Hour * hoursInDay).Add(time.Hour * time.Duration(startHours))
		} else {
			w.intervalStart = newStart
		}

		w.intervalEnd = w.intervalStart.Add(w.intervalD)
	}
}

// start starts worker, that listens to in-channel,
// collects candles from trades, handles auto-flush to file,
// when time-interval exceeds.
func (w *Worker) start() {
	cs := candles.NewStorage()

	for tr := range w.in {
		if tr.Timestamp.After(w.intervalEnd) || tr.Timestamp.Equal(w.intervalEnd) {
			w.flush(cs)
			w.incrementInterval(tr.Timestamp)
		}

		cs.AddTrade(tr, w.intervalStart)
	}

	w.flush(cs)
	close(w.out)
}

// flush flushes all data from storage to file writer.
// Does nothing if storage is empty.
func (w *Worker) flush(cs *candles.Storage) {
	if cs.Len() == 0 {
		return
	}

	c := cs.Candles()
	data := make([]string, 0, len(c))

	for i := range c {
		data = append(data, c[i].String())
	}

	chunk := strings.Join(data, "\n")
	w.out <- chunk

	cs.Clear()
}

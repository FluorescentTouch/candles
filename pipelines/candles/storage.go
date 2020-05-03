package candles

import "time"

type ticker string

// Storage stores candles for single interval.
type Storage struct {
	data map[ticker]*Candle
}

// NewStorage creates new storage.
func NewStorage() *Storage {
	return &Storage{
		data: make(map[ticker]*Candle),
	}
}

// Len returns records count in storage.
func (cs *Storage) Len() int {
	return len(cs.data)
}

// Candles returns candles for single interval.
func (cs *Storage) Candles() []Candle {
	out := make([]Candle, 0, len(cs.data))
	for _, c := range cs.data {
		out = append(out, *c)
	}

	return out
}

// Clear clears storage.
// Has to be called for each new interval.
func (cs *Storage) Clear() {
	cs.data = make(map[ticker]*Candle)
}

// AddTrade add trades to candles for single interval.
func (cs *Storage) AddTrade(trade Trade, iStart time.Time) {
	if c, ok := cs.data[trade.t]; !ok {
		cs.data[trade.t] = New(trade, iStart)
	} else {
		c.AddTrade(trade)
	}
}

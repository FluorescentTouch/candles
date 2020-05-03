package candles

import (
	"fmt"
	"time"
)

// Candle contains data about current interval deals.
type Candle struct {
	t          ticker
	startTime  time.Time
	openPrice  float64
	maxPrice   float64
	minPrice   float64
	closePrice float64
}

// Candle creates new Candle from initial Trade.
func New(trade Trade, iStart time.Time) *Candle {
	return &Candle{
		t:          trade.t,
		startTime:  iStart,
		openPrice:  trade.price,
		maxPrice:   trade.price,
		minPrice:   trade.price,
		closePrice: trade.price,
	}
}

// AddTrade adds Trade to Candle.
func (c *Candle) AddTrade(trade Trade) {
	if trade.price > c.maxPrice {
		c.maxPrice = trade.price
	}

	if trade.price < c.minPrice {
		c.minPrice = trade.price
	}

	c.closePrice = trade.price
}

// String returns string values of Candle.
func (c *Candle) String() string {
	return fmt.Sprintf("%s,%s,%f,%f,%f,%f",
		c.t,
		c.startTime.Format(time.RFC3339),
		c.openPrice,
		c.maxPrice,
		c.minPrice,
		c.closePrice,
	)
}

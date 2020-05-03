package candles

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const tradeDataLen = 4

var (
	ErrInvalidTicker = errors.New("invalid ticker provided")
	ErrInvalidValue  = errors.New("invalid value provided")
	ErrInvalidPrice  = errors.New("invalid price")
	ErrInvalidCount  = errors.New("invalid count")
	ErrInvalidTime   = errors.New("invalid timestamp")
)

// Trade contains data about trade deal.
type Trade struct {
	t         ticker
	price     float64
	count     int
	Timestamp time.Time
}

// MustTradeFromString parse a trade from a string.
// Panics if string is invalid for parsing.
func MustTradeFromString(s string) Trade {
	tr, err := TradeFromString(s)
	if err != nil {
		panic(err)
	}

	return tr
}

// TradeFromString parse Trade from string.
func TradeFromString(s string) (Trade, error) {
	values := strings.Split(strings.TrimSpace(s), ",")
	if len(values) != tradeDataLen {
		return Trade{}, ErrInvalidValue
	}

	t := values[0]
	if len(t) == 0 {
		return Trade{}, ErrInvalidTicker
	}

	price, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return Trade{}, ErrInvalidPrice
	}

	count, err := strconv.Atoi(values[2])
	if err != nil {
		return Trade{}, ErrInvalidCount
	}

	timestamp, err := time.Parse(
		"2006-01-02 15:04:05.999999",
		values[3],
	)
	if err != nil {
		return Trade{}, ErrInvalidTime
	}

	return Trade{
		t:         ticker(t),
		price:     price,
		count:     count,
		Timestamp: timestamp,
	}, nil
}

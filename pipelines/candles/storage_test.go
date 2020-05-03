package candles_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/candles/pipelines/candles"
)

func TestStorage_Clear(t *testing.T) {
	defaultTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	type args struct {
		trades []candles.Trade
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success, empty storage",
			args: args{
				trades: nil,
			},
		},
		{
			name: "success, single value",
			args: args{
				trades: []candles.Trade{
					candles.MustTradeFromString("TICKER,200.000000,10,2019-01-30 06:59:45.000249"),
				},
			},
		},
		{
			name: "success, multiple values",
			args: args{
				trades: []candles.Trade{
					candles.MustTradeFromString("TICKER_ONE,200.000000,10,2019-01-30 06:59:45.000249"),
					candles.MustTradeFromString("TICKER_TWO,100.000000,10,2019-01-30 06:59:45.000249"),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cs := candles.NewStorage()
			for _, t := range test.args.trades {
				cs.AddTrade(t, defaultTime)
			}
			cs.Clear()
			c := cs.Candles()
			assert.Equal(t, 0, len(c))
		})
	}
}

func TestStorage_Len(t *testing.T) {
	defaultTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	type args struct {
		trades []candles.Trade
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "success, empty storage",
			args: args{
				trades: nil,
			},
			want: 0,
		},
		{
			name: "success, single value",
			args: args{
				trades: []candles.Trade{
					candles.MustTradeFromString("TICKER,200.000000,10,2019-01-30 06:59:45.000249"),
				},
			},
			want: 1,
		},
		{
			name: "success, multiple values",
			args: args{
				trades: []candles.Trade{
					candles.MustTradeFromString("TICKER_ONE,200.000000,10,2019-01-30 06:59:45.000249"),
					candles.MustTradeFromString("TICKER_TWO,100.000000,10,2019-01-30 06:59:45.000249"),
				},
			},
			want: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cs := candles.NewStorage()
			for _, t := range test.args.trades {
				cs.AddTrade(t, defaultTime)
			}
			assert.Equal(t, test.want, cs.Len())
		})
	}
}

package candles_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/candles/pipelines/candles"
)

func TestCandle_String(t *testing.T) {
	defaultTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	type args struct {
		c *candles.Candle
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty candle",
			args: args{c: &candles.Candle{}},
			want: ",0001-01-01T00:00:00Z,0.000000,0.000000,0.000000,0.000000",
		},
		{
			name: "non-empty candle",
			args: args{c: candles.New(
				candles.MustTradeFromString("TICKER,213.8,100,2019-01-30 06:59:45.000249"),
				defaultTime,
			)},
			want: "TICKER,2006-01-02T15:04:05Z,213.800000,213.800000,213.800000,213.800000",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.args.c.String()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestCandle_AddTrade(t *testing.T) {
	defaultTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	type args struct {
		c *candles.Candle
		t candles.Trade
	}

	tests := []struct {
		name       string
		args       args
		wantString string
	}{
		{
			name: "empty candle, empty trade",
			args: args{
				c: &candles.Candle{},
				t: candles.Trade{},
			},
			wantString: ",0001-01-01T00:00:00Z,0.000000,0.000000,0.000000,0.000000",
		},
		{
			name: "non-empty values, max price",
			args: args{
				c: candles.New(
					candles.MustTradeFromString("TICKER,111.0,100,2019-01-30 06:59:45.000249"),
					defaultTime,
				),
				t: candles.MustTradeFromString("TICKER,222.0,60,2019-01-30 06:59:45.000249"),
			},
			wantString: "TICKER,2006-01-02T15:04:05Z,111.000000,222.000000,111.000000,222.000000",
		},
		{
			name: "non-empty values, min price",
			args: args{
				c: candles.New(
					candles.MustTradeFromString("TICKER,222.0,100,2019-01-30 06:59:45.000249"),
					defaultTime,
				),
				t: candles.MustTradeFromString("TICKER,111.0,60,2019-01-30 06:59:45.000249"),
			},
			wantString: "TICKER,2006-01-02T15:04:05Z,222.000000,222.000000,111.000000,111.000000",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := test.args.c
			c.AddTrade(test.args.t)
			assert.Equal(t, test.wantString, c.String())
		})
	}
}

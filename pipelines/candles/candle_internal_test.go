package candles

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCandle_Internal_AddTrade(t *testing.T) {
	defaultTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	type args struct {
		c *Candle
		t Trade
	}

	tests := []struct {
		name string
		args args
		want *Candle
	}{
		{
			name: "empty candle, empty trade",
			args: args{
				c: &Candle{},
				t: Trade{},
			},
			want: &Candle{},
		},
		{
			name: "non-empty values, max price",
			args: args{
				c: &Candle{
					t:          ticker("TICKER"),
					startTime:  defaultTime,
					openPrice:  100.0,
					maxPrice:   200.0,
					minPrice:   50.0,
					closePrice: 150.0,
				},
				t: MustTradeFromString("TICKER,300.0,60,2019-01-30 06:59:45.000249"),
			},
			want: &Candle{
				t:          ticker("TICKER"),
				startTime:  defaultTime,
				openPrice:  100.0,
				maxPrice:   300.0,
				minPrice:   50.0,
				closePrice: 300.0,
			},
		},
		{
			name: "non-empty values, min price",
			args: args{
				c: &Candle{
					t:          ticker("TICKER"),
					startTime:  defaultTime,
					openPrice:  100.0,
					maxPrice:   200.0,
					minPrice:   50.0,
					closePrice: 150.0,
				},
				t: MustTradeFromString("TICKER,25.0,60,2019-01-30 06:59:45.000249"),
			},
			want: &Candle{
				t:          ticker("TICKER"),
				startTime:  defaultTime,
				openPrice:  100.0,
				maxPrice:   200.0,
				minPrice:   25.0,
				closePrice: 25.0,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := test.args.c
			c.AddTrade(test.args.t)
			assert.Equal(t, test.want, c)
		})
	}
}

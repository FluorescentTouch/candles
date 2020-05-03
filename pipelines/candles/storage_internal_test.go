package candles

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStorage_Internal_Candles(t *testing.T) {
	defaultTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	type args struct {
		cs Storage
	}

	tests := []struct {
		name string
		args args
		want map[Candle]struct{}
	}{
		{
			name: "success, multiple values",
			args: args{cs: Storage{
				data: map[ticker]*Candle{
					ticker("TICKER_ONE"): {
						t:          ticker("TICKER_ONE"),
						startTime:  defaultTime,
						openPrice:  100.0,
						maxPrice:   100.0,
						minPrice:   100.0,
						closePrice: 100.0,
					},
					ticker("TICKER_TWO"): {
						t:          ticker("TICKER_TWO"),
						startTime:  defaultTime,
						openPrice:  200.0,
						maxPrice:   200.0,
						minPrice:   200.0,
						closePrice: 200.0,
					},
				},
			}},
			want: map[Candle]struct{}{
				{
					t:          ticker("TICKER_ONE"),
					startTime:  defaultTime,
					openPrice:  100.0,
					maxPrice:   100.0,
					minPrice:   100.0,
					closePrice: 100.0,
				}: {},
				{
					t:          ticker("TICKER_TWO"),
					startTime:  defaultTime,
					openPrice:  200.0,
					maxPrice:   200.0,
					minPrice:   200.0,
					closePrice: 200.0,
				}: {},
			},
		},
		{
			name: "empty storage",
			args: args{cs: Storage{
				data: map[ticker]*Candle{},
			}},
			want: map[Candle]struct{}{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.args.cs.Candles()

			out := make(map[Candle]struct{})
			for _, c := range got {
				out[c] = struct{}{}
			}
			assert.Equal(t, test.want, out)
		})
	}
}

func TestStorage_Internal_AddTrade(t *testing.T) {
	defaultTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	intervalStartTime, _ := time.Parse(time.RFC3339, "2006-01-02T10:00:00Z")

	type args struct {
		cs     *Storage
		t      Trade
		iStart time.Time
	}

	tests := []struct {
		name string
		args args
		want *Storage
	}{
		{
			name: "success, first value",
			args: args{
				cs: &Storage{
					data: map[ticker]*Candle{},
				},
				t: Trade{
					t:         ticker("TICKER"),
					price:     100.0,
					count:     10,
					Timestamp: defaultTime,
				},
				iStart: intervalStartTime,
			},
			want: &Storage{
				data: map[ticker]*Candle{
					ticker("TICKER"): {
						t:          ticker("TICKER"),
						startTime:  intervalStartTime,
						openPrice:  100.0,
						maxPrice:   100.0,
						minPrice:   100.0,
						closePrice: 100.0,
					},
				},
			},
		},
		{
			name: "success, new ticker added",
			args: args{
				cs: &Storage{
					data: map[ticker]*Candle{
						ticker("TICKER_ONE"): {
							t: ticker("TICKER_ONE"),
						},
					},
				},
				t: Trade{
					t:         ticker("TICKER_TWO"),
					price:     100.0,
					count:     10,
					Timestamp: defaultTime,
				},
				iStart: intervalStartTime,
			},
			want: &Storage{
				data: map[ticker]*Candle{
					ticker("TICKER_ONE"): {
						t: ticker("TICKER_ONE"),
					},
					ticker("TICKER_TWO"): {
						t:          ticker("TICKER_TWO"),
						startTime:  intervalStartTime,
						openPrice:  100.0,
						maxPrice:   100.0,
						minPrice:   100.0,
						closePrice: 100.0,
					},
				},
			},
		},
		{
			name: "success, add to existing ticker",
			args: args{
				cs: &Storage{
					data: map[ticker]*Candle{
						ticker("TICKER"): {
							t:          ticker("TICKER"),
							startTime:  intervalStartTime,
							openPrice:  100.0,
							maxPrice:   100.0,
							minPrice:   100.0,
							closePrice: 100.0,
						},
					},
				},
				t: Trade{
					t:         ticker("TICKER"),
					price:     200.0,
					count:     10,
					Timestamp: defaultTime,
				},
				iStart: intervalStartTime,
			},
			want: &Storage{
				data: map[ticker]*Candle{
					ticker("TICKER"): {
						t:          ticker("TICKER"),
						startTime:  intervalStartTime,
						openPrice:  100.0,
						maxPrice:   200.0,
						minPrice:   100.0,
						closePrice: 200.0,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.args.cs.AddTrade(test.args.t, test.args.iStart)
			assert.Equal(t, test.want, test.args.cs)
		})
	}
}

func TestStorage_Internal_Clear(t *testing.T) {
	defaultTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	type args struct {
		cs *Storage
	}

	tests := []struct {
		name string
		args args
		want *Storage
	}{
		{
			name: "success, multiple values",
			args: args{cs: &Storage{
				data: map[ticker]*Candle{
					ticker("TICKER_ONE"): {
						t:          ticker("TICKER_ONE"),
						startTime:  defaultTime,
						openPrice:  100.0,
						maxPrice:   100.0,
						minPrice:   100.0,
						closePrice: 100.0,
					},
					ticker("TICKER_TWO"): {
						t:          ticker("TICKER_TWO"),
						startTime:  defaultTime,
						openPrice:  200.0,
						maxPrice:   200.0,
						minPrice:   200.0,
						closePrice: 200.0,
					},
				},
			}},
			want: &Storage{
				data: map[ticker]*Candle{},
			},
		},
		{
			name: "empty storage",
			args: args{cs: &Storage{
				data: map[ticker]*Candle{},
			}},
			want: &Storage{
				data: map[ticker]*Candle{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.args.cs.Clear()
			assert.Equal(t, test.want, test.args.cs)
		})
	}
}

func TestStorage_Internal_Len(t *testing.T) {
	defaultTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	type args struct {
		cs *Storage
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success, multiple values",
			args: args{cs: &Storage{
				data: map[ticker]*Candle{
					ticker("TICKER_ONE"): {
						t:          ticker("TICKER_ONE"),
						startTime:  defaultTime,
						openPrice:  100.0,
						maxPrice:   100.0,
						minPrice:   100.0,
						closePrice: 100.0,
					},
					ticker("TICKER_TWO"): {
						t:          ticker("TICKER_TWO"),
						startTime:  defaultTime,
						openPrice:  200.0,
						maxPrice:   200.0,
						minPrice:   200.0,
						closePrice: 200.0,
					},
				},
			}},
		},
		{
			name: "empty storage",
			args: args{cs: &Storage{
				data: map[ticker]*Candle{},
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, len(test.args.cs.data), test.args.cs.Len())
		})
	}
}

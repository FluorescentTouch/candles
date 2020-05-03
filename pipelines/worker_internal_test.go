package pipelines

import (
	"strings"
	"testing"
	"time"

	"github.com/candles/pipelines/candles"

	"github.com/stretchr/testify/assert"
)

func mustParseTime(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05.999999", s)
	if err != nil {
		panic(err)
	}

	return t
}

func TestWorker_Internal_flush(t *testing.T) {
	defaultTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	type args struct {
		w      *Worker
		trades []candles.Trade
	}

	tests := []struct {
		name       string
		args       args
		wantOutput bool
		want       map[string]bool
	}{
		{
			name: "success, multiple values in storage",
			args: args{
				w: &Worker{
					out: make(chan string, 1),
				},
				trades: []candles.Trade{
					candles.MustTradeFromString("TICKER_ONE,200.000000,10,2019-01-30 06:59:45.000249"),
					candles.MustTradeFromString("TICKER_TWO,100.000000,10,2019-01-30 06:59:45.000249"),
				},
			},
			wantOutput: true,
			want: map[string]bool{
				"TICKER_ONE,2006-01-02T15:04:05Z,200.000000,200.000000,200.000000,200.000000": true,
				"TICKER_TWO,2006-01-02T15:04:05Z,100.000000,100.000000,100.000000,100.000000": true,
			},
		},
		{
			name: "empty storage",
			args: args{
				w: &Worker{
					out: make(chan string, 1),
				},
			},
			wantOutput: false,
			want:       map[string]bool{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cs := candles.NewStorage()
			for _, tr := range test.args.trades {
				cs.AddTrade(tr, defaultTime)
			}
			test.args.w.flush(cs)
			defer close(test.args.w.out)
			var (
				output    string
				outExists bool
			)
			select {
			case output = <-test.args.w.out:
				outExists = true
			default:
			}
			outCheck := make(map[string]bool)
			if output != "" {
				for _, s := range strings.Split(output, "\n") {
					outCheck[s] = true
				}
			}
			assert.Equal(t, test.wantOutput, outExists)
			assert.Equal(t, test.want, outCheck)
		})
	}
}

func TestWorker_Internal_start(t *testing.T) {
	defaultTime := mustParseTime("2019-01-30 11:00:00.000000")
	tradeOne := candles.MustTradeFromString("TICKER_ONE,200.000000,10,2019-01-30 11:30:45.000000")
	tradeTwo := candles.MustTradeFromString("TICKER_TWO,100.000000,10,2019-01-30 11:59:45.000000")

	type args struct {
		w      *Worker
		trades []candles.Trade
	}

	tests := []struct {
		name       string
		args       args
		wantOutput bool
		want       map[string]bool
	}{
		{
			name: "success, multiple values in storage",
			args: args{
				w: &Worker{
					out: make(chan string),
					in:  make(chan candles.Trade),

					intervalStart: defaultTime,
					intervalEnd:   defaultTime.Add(time.Hour),
				},
				trades: []candles.Trade{
					tradeOne,
					tradeTwo,
				},
			},
			wantOutput: true,
			want: map[string]bool{
				"TICKER_ONE,2019-01-30T11:00:00Z,200.000000,200.000000,200.000000,200.000000": true,
				"TICKER_TWO,2019-01-30T11:00:00Z,100.000000,100.000000,100.000000,100.000000": true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			go test.args.w.start()
			for _, tr := range test.args.trades {
				test.args.w.in <- tr
			}
			close(test.args.w.in)
			var (
				output    string
				outExists bool
			)
			select {
			case output = <-test.args.w.out:
				outExists = true
			case <-time.NewTicker(time.Second * 2).C:
			}
			outCheck := make(map[string]bool)
			if output != "" {
				for _, s := range strings.Split(output, "\n") {
					outCheck[s] = true
				}
			}
			assert.Equal(t, test.wantOutput, outExists)
			assert.Equal(t, test.want, outCheck)
		})
	}
}

func TestWorker_Internal_incrementInterval(t *testing.T) {
	defaultTime := mustParseTime("2019-01-30 11:00:00.000000")
	defaultInterval := 5
	defaultIDuration := time.Minute * time.Duration(defaultInterval)

	type args struct {
		w  *Worker
		tr candles.Trade
	}

	tests := []struct {
		name      string
		args      args
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name: "success, first trade",
			args: args{
				w: &Worker{
					interval:  defaultInterval,
					intervalD: defaultIDuration,
				},
				tr: candles.MustTradeFromString("TICKER_ONE,200.000000,10,2019-01-30 11:02:00.000000"),
			},
			wantStart: mustParseTime("2019-01-30 11:00:00.000000"),
			wantEnd:   mustParseTime("2019-01-30 11:05:00.000000"),
		},
		{
			name: "success, normal increment",
			args: args{
				w: &Worker{
					interval:  defaultInterval,
					intervalD: defaultIDuration,

					intervalStart: defaultTime,
					intervalEnd:   defaultTime.Add(defaultIDuration),
				},
				tr: candles.MustTradeFromString("TICKER_ONE,200.000000,10,2019-01-31 11:02:00.000000"),
			},
			wantStart: mustParseTime("2019-01-31 11:00:00.000000"),
			wantEnd:   mustParseTime("2019-01-31 11:05:00.000000"),
		},
		{
			name: "success, increment for not-working period",
			args: args{
				w: &Worker{
					interval:      defaultInterval,
					intervalD:     defaultIDuration,
					intervalStart: mustParseTime("2019-01-30 04:00:00.000000"),
					intervalEnd:   mustParseTime("2019-01-30 04:05:00.000000"),
				},
				tr: candles.MustTradeFromString("TICKER_ONE,200.000000,10,2019-01-31 11:00:00.000000"),
			},
			wantStart: mustParseTime("2019-01-31 11:00:00.000000"),
			wantEnd:   mustParseTime("2019-01-31 11:05:00.000000"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.args.w.incrementInterval(test.args.tr.Timestamp)
			assert.Equal(t, test.wantStart, test.args.w.intervalStart)
			assert.Equal(t, test.wantEnd, test.args.w.intervalEnd)
		})
	}
}

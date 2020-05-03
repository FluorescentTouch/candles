package candles_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/candles/pipelines/candles"
)

func TestTradeFromString(t *testing.T) {
	type args struct {
		s string
	}

	mustParseTime := func(s string) time.Time {
		t, err := time.Parse("2006-01-02 15:04:05.999999", s)
		if err != nil {
			panic(err)
		}

		return t
	}

	tests := []struct {
		name          string
		args          args
		wantTimeStamp time.Time
		wantErr       error
	}{
		{
			name:          "success",
			args:          args{s: "TICKER,213.8,10,2019-01-30 06:59:45.000249"},
			wantTimeStamp: mustParseTime("2019-01-30 06:59:45.000249"),
			wantErr:       nil,
		},
		{
			name:          "num of elements exceeded",
			args:          args{s: "TICKER,213.8,10,2019-01-30 06:59:45.000249,new-info"},
			wantTimeStamp: time.Time{},
			wantErr:       candles.ErrInvalidValue,
		},
		{
			name:          "invalid ticker",
			args:          args{s: ",213.8,10,2019-01-30 06:59:45.000249"},
			wantTimeStamp: time.Time{},
			wantErr:       candles.ErrInvalidTicker,
		},
		{
			name:          "invalid price",
			args:          args{s: "TICKER,ohe hundred dollars,10,2019-01-30 06:59:45.000249"},
			wantTimeStamp: time.Time{},
			wantErr:       candles.ErrInvalidPrice,
		},
		{
			name:          "invalid count",
			args:          args{s: "TICKER,213.8,five,2019-01-30 06:59:45.000249"},
			wantTimeStamp: time.Time{},
			wantErr:       candles.ErrInvalidCount,
		},
		{
			name:          "invalid timestamp",
			args:          args{s: "TICKER,213.8,10,three-o-clock"},
			wantTimeStamp: time.Time{},
			wantErr:       candles.ErrInvalidTime,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := candles.TradeFromString(test.args.s)
			assert.Equal(t, test.wantTimeStamp, got.Timestamp)
			assert.Equal(t, test.wantErr, err)
		})
	}
}

func TestMustTradeFromString(t *testing.T) {
	type args struct {
		s string
	}

	mustParseTime := func(s string) time.Time {
		t, err := time.Parse("2006-01-02 15:04:05.999999", s)
		if err != nil {
			panic(err)
		}

		return t
	}

	tests := []struct {
		name          string
		args          args
		wantTimeStamp time.Time
		wantErr       error
	}{
		{
			name:          "success",
			args:          args{s: "TICKER,213.8,10,2019-01-30 06:59:45.000249"},
			wantTimeStamp: mustParseTime("2019-01-30 06:59:45.000249"),
			wantErr:       nil,
		},
		{
			name:          "panic",
			args:          args{s: ",213.8,10,2019-01-30 06:59:45.000249"},
			wantTimeStamp: mustParseTime("2019-01-30 06:59:45.000249"),
			wantErr:       candles.ErrInvalidTicker,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()
				assert.Equal(t, test.wantErr, r)
			}()
			got := candles.MustTradeFromString(test.args.s)
			assert.Equal(t, test.wantTimeStamp, got.Timestamp)
		})
	}
}

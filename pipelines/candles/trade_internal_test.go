package candles

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInternalTradeFromString(t *testing.T) {
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
		name    string
		args    args
		want    Trade
		wantErr error
	}{
		{
			name: "success",
			args: args{s: "TICKER,213.8,10,2019-01-30 06:59:45.000249"},
			want: Trade{
				t:         ticker("TICKER"),
				price:     213.8,
				count:     10,
				Timestamp: mustParseTime("2019-01-30 06:59:45.000249"),
			},
			wantErr: nil,
		},
		{
			name:    "num of elements exceeded",
			args:    args{s: "TICKER,213.8,10,2019-01-30 06:59:45.000249,new-info"},
			want:    Trade{},
			wantErr: ErrInvalidValue,
		},
		{
			name:    "invalid ticker",
			args:    args{s: ",213.8,10,2019-01-30 06:59:45.000249"},
			want:    Trade{},
			wantErr: ErrInvalidTicker,
		},
		{
			name:    "invalid price",
			args:    args{s: "TICKER,ohe hundred dollars,10,2019-01-30 06:59:45.000249"},
			want:    Trade{},
			wantErr: ErrInvalidPrice,
		},
		{
			name:    "invalid count",
			args:    args{s: "TICKER,213.8,five,2019-01-30 06:59:45.000249"},
			want:    Trade{},
			wantErr: ErrInvalidCount,
		},
		{
			name:    "invalid timestamp",
			args:    args{s: "TICKER,213.8,10,three-o-clock"},
			want:    Trade{},
			wantErr: ErrInvalidTime,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := TradeFromString(test.args.s)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err)
		})
	}
}

func TestInternalMustTradeFromString(t *testing.T) {
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
		name    string
		args    args
		want    Trade
		wantErr error
	}{
		{
			name: "success",
			args: args{s: "TICKER,213.8,10,2019-01-30 06:59:45.000249"},
			want: Trade{
				t:         ticker("TICKER"),
				price:     213.8,
				count:     10,
				Timestamp: mustParseTime("2019-01-30 06:59:45.000249"),
			},
			wantErr: nil,
		},
		{
			name:    "panic",
			args:    args{s: ",213.8,10,2019-01-30 06:59:45.000249"},
			want:    Trade{},
			wantErr: ErrInvalidTicker,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()
				assert.Equal(t, test.wantErr, r)
			}()
			got := MustTradeFromString(test.args.s)
			assert.Equal(t, test.want, got)
		})
	}
}

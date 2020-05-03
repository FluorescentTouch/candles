package pipelines

import (
	"errors"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestWorker_Internal_startWriting(t *testing.T) {
	writerMock := NewWriterMock(t)
	setWriteString := func(s string, err error) {
		writerMock.On(
			"WriteString",
			s+"\n",
		).Return(err).Once()
	}

	l := logrus.New()

	type args struct {
		data []string
	}

	tests := []struct {
		name  string
		setup func(a args)
		args  args
	}{
		{
			name: "success, multiple values written",
			args: args{
				data: []string{
					"stringOne",
					"stringTwo",
				},
			},
			setup: func(a args) {
				for _, s := range a.data {
					setWriteString(s, nil)
				}
			},
		},
		{
			name: "multiple values written, error on write",
			args: args{
				data: []string{
					"stringOne",
					"stringTwo",
				},
			},
			setup: func(a args) {
				for _, s := range a.data {
					setWriteString(s, errors.New("something went wrong"))
				}
			},
		},
		{
			name: "zero values written",
			args: args{
				data: []string{},
			},
			setup: func(a args) {
				for _, s := range a.data {
					setWriteString(s, nil)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(test.args)
			defer writerMock.AssertExpectations(t)

			input := make(chan string)

			w := NewWriter(writerMock, input, l)
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go w.startWriting(wg)

			for _, s := range test.args.data {
				input <- s
			}
			close(input)
			wg.Wait()
		})
	}
}

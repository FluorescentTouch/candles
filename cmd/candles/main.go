package main

import (
	"flag"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/candles/files"
	"github.com/candles/pipelines"
)

var filepath string

const timeToWait = 5

func main() {
	flag.StringVar(&filepath, "filepath", "trades.csv", "path to files with trades")
	flag.Parse()

	logger := logrus.New()

	reader, err := files.NewReader(filepath, logger)
	if err != nil {
		logger.Errorf("can't init file reader: %v", err)
		os.Exit(1)
	}

	wrBuilder := pipelines.WriterBuilder{}
	p := pipelines.New(reader, wrBuilder, logger)

	for _, interval := range []int{5, 30, 240} {
		err = p.Add(interval)
		if err != nil {
			logger.Errorf("can't add pipeline to pipelines: %v", err)
			os.Exit(1)
		}
	}

	// init pipeline
	go p.Init()

	// send starting signal
	p.Start <- struct{}{}

	// waiting for file reading end
	<-p.FileDone

	timeout := time.NewTicker(time.Second * timeToWait)
	select {
	case <-timeout.C:
		logger.Errorf("pipeline end timeout exceeded")
		os.Exit(1)
	case <-p.Done:
		logger.Info("Successfully completed")
	}
}

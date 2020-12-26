package main

import (
	"context"
	"os"
	"sync"

	"github.com/andrewneudegg/delta/cmd/common"
	"github.com/andrewneudegg/delta/pkg/sink"
	log "github.com/sirupsen/logrus"
)

func main() {
	common.LoggingInit()
	probes := common.ProbesInit()
	common.TelemetryInit()

	wg := sync.WaitGroup{}
	wg.Add(1)
	log.Info("starting application")

	probes.ReadyNow()
	app()

	wg.Wait()
	log.Warn("exiting application")
}

func app() {
	mq := make(chan *sink.SunkMessage)

	sinkServer, err := sink.NewHTTPSinkServer(&sink.HTTPSinkServerConfiguration{
		ServerConfiguration: sink.ServerConfiguration{
			ToChan: mq,
		},
		ListenAddr:  ":8080",
		MaxBodySize: 2097152, // two Mebibytes
	})

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// just make the chan work...
	go func() {
		for {
			<-mq
		}
	}()

	go func() {
		err := sinkServer.Serve(context.TODO())
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}()
}

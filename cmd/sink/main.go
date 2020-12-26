package main

import (
	"context"
	"os"
	"sync"

	"github.com/andrewneudegg/delta/pkg/probes"
	"github.com/andrewneudegg/delta/pkg/sink"
	log "github.com/sirupsen/logrus"
)

func init() {
	// log.SetReportCaller(true)
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	// API that listens to anything.
	wg := sync.WaitGroup{}
	wg.Add(1)

	log.Info("starting application")

	// Probes
	probes := probes.ProbeServer{
		Port: 8082,
	}
	go probes.StartProbeServer()

	probes.AliveNow()
	probes.ReadyNow()

	mq := make(chan *sink.SunkMessage)

	sinkServer, err := sink.NewHTTPSinkServer(&sink.HTTPSinkServerConfiguration{
		ServerConfiguration: sink.ServerConfiguration{
			ToChan: mq,
		},
		ListenAddr: ":8080",
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

	wg.Wait()
	log.Warn("exiting application")
}

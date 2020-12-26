package main

import (
	"sync"

	"github.com/andrewneudegg/delta/cmd/common"
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
	return
}

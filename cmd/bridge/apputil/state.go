package apputil

import (
	"context"
	"sync"

	"github.com/andrewneudegg/delta/pkg/probes"
	"github.com/andrewneudegg/delta/pkg/telemetry"
	log "github.com/sirupsen/logrus"
)

// AppState contains the information for the app to run.
type AppState struct {
	Probes          *probes.ProbeServer
	Logger          *log.Logger
	TelemetryServer *telemetry.PrometheusServer
	wg              sync.WaitGroup
}

// Block the application indefinitely on the given context.
func (a *AppState) Block(ctx context.Context) {
	<-ctx.Done()
}

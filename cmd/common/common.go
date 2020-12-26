// Package common for all applications
package common

import (
	"context"
	"os"

	"github.com/andrewneudegg/delta/pkg/probes"
	"github.com/andrewneudegg/delta/pkg/telemetry"

	log "github.com/sirupsen/logrus"
)

// LoggingInit will configure logrus to ensure consistency accross all applications.
func LoggingInit() {
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

// ProbesInit will create and activate the liveness probe server.
func ProbesInit() probes.ProbeServer {
	probes := probes.ProbeServer{
		ListenAddr: ":8082",
	}
	go probes.StartProbeServer()
	probes.AliveNow()
	return probes
}

// TelemetryInit will create and activate the Prometheus telemetry server.
func TelemetryInit() telemetry.PrometheusServer {
	// --------- Telemetry
	prometheusServer := telemetry.PrometheusServer{
		ListenAddr: ":8081",
		Route:      "",
	}
	go prometheusServer.Serve(context.TODO())
	return prometheusServer
}

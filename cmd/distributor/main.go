package main

import (
	"context"
	"os"

	"github.com/andrewneudegg/delta/cmd/distributor/apputil"
	"github.com/andrewneudegg/delta/cmd/distributor/subcmd/naive"
	"github.com/andrewneudegg/delta/pkg/probes"
	"github.com/andrewneudegg/delta/pkg/telemetry"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	app()
}

func configureLogger(verbose bool) *log.Logger {
	logger := log.New()

	logger.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	if verbose {
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}

	return logger
}

func configureProbes(probesEnabled bool) *probes.ProbeServer {
	probes := probes.ProbeServer{
		ListenAddr: ":8082",
	}

	if probesEnabled {
		go probes.StartProbeServer()
	}

	probes.AliveNow()
	return &probes
}

func configureTelemetryServer(telemetryEnabled bool) *telemetry.PrometheusServer {
	prometheusServer := telemetry.PrometheusServer{
		ListenAddr: ":8081",
		Route:      "/metrics",
	}

	if telemetryEnabled {
		go prometheusServer.Serve(context.TODO())
	}

	return &prometheusServer
}

func app() error {
	var verboseMode bool
	var telemetryEnabled bool
	var probesEnabled bool

	appState := apputil.AppState{
		Probes:          nil,
		Logger:          nil,
		TelemetryServer: nil,
	}

	var rootCmd = &cobra.Command{
		Use: "distributor",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			appState.Logger = configureLogger(verboseMode)
			appState.Probes = configureProbes(probesEnabled)
			appState.TelemetryServer = configureTelemetryServer(telemetryEnabled)
			return nil
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verboseMode, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&telemetryEnabled, "telemetry", "", false, "telemetry server")
	rootCmd.PersistentFlags().BoolVarP(&probesEnabled, "probes", "", false, "liveness / readiness probes")

	rootCmd.AddCommand(naive.Cmd(&appState))
	return rootCmd.Execute()
}

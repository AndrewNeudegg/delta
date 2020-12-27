// Package simulator will generate arbitrary events for integration testing
package simulator

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/andrewneudegg/delta/cmd/bridge/apputil"
)

type simulatorOpts struct {
	Interval   time.Duration
	ListenAddr string
}

// Cmd demonstrates how to configure a new subcommand.
func Cmd(appState *apputil.AppState) *cobra.Command {
	var generationInterval time.Duration
	var listenAddr string

	cmd := &cobra.Command{
		Use:   "simulator",
		Short: "simulate arbitrary events for testing",
		Long:  `simulate an arbitrary number of events for testing`,
		Run: func(cmd *cobra.Command, args []string) {
			appState.Logger.Debug("starting example cmd")
			appState.Block(context.TODO())
		},
	}

	cmd.Flags().DurationVarP(&generationInterval, "interval", "i", time.Second, "the generation interval for events")
	cmd.Flags().StringVarP(&listenAddr, "listen-address", "l", ":8080", "specify the address to listen for coordination instructions")
	return cmd
}

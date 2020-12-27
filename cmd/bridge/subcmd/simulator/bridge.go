// Package simulator will generate arbitrary events for integration testing
package simulator

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/andrewneudegg/delta/cmd/bridge/apputil"
)

// Cmd demonstrates how to configure a new subcommand.
func Cmd(appState *apputil.AppState) *cobra.Command {
	return &cobra.Command{
		Use:   "simulator",
		Short: "simulate arbitrary events for testing",
		Long:  `simulate an arbitary number of events for testing`,
		Run: func(cmd *cobra.Command, args []string) {
			appState.Logger.Debug("starting example cmd")
			appState.Block(context.TODO())
		},
	}
}

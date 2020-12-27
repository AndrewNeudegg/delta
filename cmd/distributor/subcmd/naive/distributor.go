package naive

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/andrewneudegg/delta/cmd/distributor/apputil"
)

// Cmd demonstrates how to configure a new subcommand.
func Cmd(appState *apputil.AppState) *cobra.Command {
	return &cobra.Command{
		Use:   "naive",
		Short: "naively distribute events to a specific target/s",
		Long:  `naively distribute events to a specific target/s without any leader election`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appState.Logger.Debug("starting example cmd")
			appState.Block(context.TODO())
		},
	}
}

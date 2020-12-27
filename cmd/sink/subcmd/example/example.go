package example

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/andrewneudegg/delta/cmd/bridge/apputil"
)

// Cmd demonstrates how to configure a new subcommand.
func Cmd(appState *apputil.AppState) *cobra.Command {
	return &cobra.Command{
		Use:   "example [string to print]",
		Short: "example anything to the screen",
		Long: `example is for printing anything back to the screen.
For many years people have printed back to the screen.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appState.Logger.Debug("starting example cmd")
			appState.Block(context.TODO())
		},
	}
}

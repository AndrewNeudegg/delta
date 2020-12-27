package naive

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/andrewneudegg/delta/cmd/distributor/apputil"
)

// Cmd demonstrates how to configure a new subcommand.
func Cmd(appState *apputil.AppState) *cobra.Command {
	var eventSources []string
	var targetServices []string

	cmd := &cobra.Command{
		Use:   "naive",
		Short: "naively distribute events to a specific target/s",
		Long:  `naively distribute events to a specific target/s without any leader election`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(eventSources) == 0 {
				return fmt.Errorf("must specify at least one event source")
			}

			if len(targetServices) == 0 {
				return fmt.Errorf("must specify at least one target service")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			appState.Logger.Debug("starting example cmd")
			appState.Block(context.TODO())
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&eventSources, "event-sources", "e", []string{}, "sources of events")
	cmd.Flags().StringSliceVarP(&targetServices, "target-services", "t", []string{}, "the target services to send events to")

	return cmd
}

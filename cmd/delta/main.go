// Package main is the entrypoint of the delta application.
package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/andrewneudegg/delta/cmd/delta/subcmd/serve"
)

func configureLogger(verbose bool) {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	var verboseMode bool

	var rootCmd = &cobra.Command{
		Use: "delta",
		Long: `Delta: An easy to understand and pluggable eventing system.
For more information and up-to-date documentation take a look at http://github.com/AndrewNeudegg/delta.
		`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			configureLogger(verboseMode)
			return nil
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verboseMode, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(serve.Cmd())

	rootCmd.Execute()
}

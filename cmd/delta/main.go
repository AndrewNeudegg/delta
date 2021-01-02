// Package main is the entrypoint of the delta application.
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/andrewneudegg/delta/cmd/delta/subcmd/serve"
)

func configureLogger(verbose bool) {

	packageName := "/delta/"

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	if verbose {
		log.SetReportCaller(true)
		log.StandardLogger().SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				// s := strings.Split(f.Function, ".")
				// funcName := s[len(s)-1]
				relativeFPath := strings.SplitAfterN(f.File, packageName, 2)[1]
				return "", fmt.Sprintf(" %s:%d", relativeFPath, f.Line)
				// return funcName, fmt.Sprintf(" %s:%d", relativeFPath, f.Line)
			},

			DisableColors: false,
			FullTimestamp: true,
		})
		log.SetLevel(log.DebugLevel)
	} else {
		log.StandardLogger().SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
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

	err := rootCmd.Execute()
	if err != nil {
		log.Error(err)
	}
	log.Warn("application exiting")
}

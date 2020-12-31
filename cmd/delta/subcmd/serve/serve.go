package serve

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/andrewneudegg/delta/pkg/configuration"
	"github.com/andrewneudegg/delta/pkg/pipeline"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// Cmd demonstrates how to configure a new subcommand.
func Cmd() *cobra.Command {
	var configurationPath string
	var config configuration.Container

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "serve event routing as described by the given configuration",
		Long:  `serve will orchestrate whatever recipe has been described by your configuration.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {

			configurationBytes := []byte{}
			if configurationPath == "-" {
				scanner := bufio.NewScanner(os.Stdin)
				if !scanner.Scan() {
					return errors.Wrap(scanner.Err(), "failed to read configuration")
				}
				configurationBytes = scanner.Bytes()
			} else {
				data, err := ioutil.ReadFile(configurationPath)
				if err != nil {
					return err
				}
				configurationBytes = data
			}

			cLoader := configuration.RawConfig{
				ConfigData: configurationBytes,
			}

			c, err := cLoader.Load()
			if err != nil {
				return errors.Wrap(err, "failed to load configuration from bytes")
			}

			config = c

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = config.ApplicationSettings

			log.Info("starting serve")
			p, err := pipeline.BuildPipeline(config)
			if err != nil {
				return errors.Wrap(err, "failed to build application pipeline")
			}
			p.Await()
			log.Info("stopping serve")

			return nil
		},
	}

	cmd.Flags().StringVarP(&configurationPath, "config", "c", "", "the application configuration, if using stdin specify '-'.")

	return cmd
}

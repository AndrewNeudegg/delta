package serve

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/andrewneudegg/delta/pkg/configuration"
	"github.com/andrewneudegg/delta/pkg/pipelines"
	pipelineBuilder "github.com/andrewneudegg/delta/pkg/pipelines/builder"
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
			
			var configurationBytes []byte
			if configurationPath == "-" {
				scanner := bufio.NewScanner(os.Stdin)
				if !scanner.Scan() {
					return errors.Wrap(scanner.Err(), "failed to read configuration")
				}
				configurationBytes = scanner.Bytes()
			} else {
				if configurationPath == "" {
					return fmt.Errorf("no configuration file supplied")
				}

				data, err := ioutil.ReadFile(configurationPath)
				if err != nil {
					return err
				}
				configurationBytes = data
			}

			c, err := configuration.FromBytes(configurationBytes)
			if err != nil {
				return errors.Wrap(err, "failed to load configuration from bytes")
			}

			config = c

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			wg := sync.WaitGroup{}

			// Build all pipelines.
			for _, pipe := range config.Pipeline {
				wg.Add(1)
				p, err := pipelines.BuildPipeline(pipe.ID, pipe.Config, pipelineBuilder.PipelineMapping())
				if err != nil {
					return errors.Wrapf(err, "failed to build pipeline '%s'", pipe.ID)
				}

				go func() {
					err := p.Do(context.Background())

					if err != nil {
						log.Error(err)
					}

					wg.Done()
				}()
			}

			wg.Wait()
			return nil
		},
	}

	cmd.Flags().StringVarP(&configurationPath, "config", "c", "", "the application configuration, if using stdin specify '-'.")

	return cmd
}

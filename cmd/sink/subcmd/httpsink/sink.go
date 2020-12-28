package httpsink

import (
	"context"
	"os"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/sink"
	"github.com/andrewneudegg/delta/pkg/utils"

	"github.com/spf13/cobra"

	"github.com/andrewneudegg/delta/cmd/sink/apputil"
)

// httpSinkOpts op
type httpSinkOpts struct {
	ListenAddr  string
	MaxBodySize int64

	appState *apputil.AppState
}

// Cmd demonstrates how to configure a new subcommand.
func Cmd(appState *apputil.AppState) *cobra.Command {
	var listenAddr string
	var maxBodySizeStr string

	cmd := &cobra.Command{
		Use:   "http",
		Short: "http sink for arbitrary event data",
		Long:  `http will capture any http post data to any route on the given port`,
		RunE: func(cmd *cobra.Command, args []string) error {
			appState.Logger.Debug("starting http sink")

			maxBodySize, err := utils.LabelledBytes2Int64(maxBodySizeStr)
			if err != nil {
				return err
			}

			go httpSink(&httpSinkOpts{
				ListenAddr:  listenAddr,
				MaxBodySize: maxBodySize,
				appState:    appState,
			})

			appState.Block(context.TODO())
			return nil
		},
	}

	cmd.Flags().StringVarP(&listenAddr, "listenaddr", "", ":8080", "Address to listen to events on.")
	cmd.Flags().StringVarP(&maxBodySizeStr, "maxbodysize", "", "1Mb", "Max number of bytes to accept in payload.")

	return cmd
}

func httpSink(opts *httpSinkOpts) {
	log := opts.appState.Logger // continence until I am sure this is what I want to do.

	log.Debugf("Running http sink with config: %+v", opts)

	log.Info("here")
	mq := make(chan events.Event)

	sinkServer, err := sink.NewHTTPSinkServer(&sink.HTTPSinkServerConfiguration{
		ServerConfiguration: sink.ServerConfiguration{
			ToChan: mq,
		},
		ListenAddr:  ":8080",
		MaxBodySize: 2097152,
	})

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// just make the chan work...
	go func() {
		for {
			<-mq
		}
	}()

	go func() {
		err := sinkServer.Serve(context.TODO())
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}()
}

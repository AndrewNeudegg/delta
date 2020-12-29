// Package simulator will generate arbitrary events for integration testing
package simulator

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/andrewneudegg/delta/cmd/bridge/apputil"
)

type simulatorOpts struct {
	Interval time.Duration

	appState *apputil.AppState
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

		},
	}

	cmd.Flags().DurationVarP(&generationInterval, "interval", "i", time.Second, "the generation interval for events")
	cmd.Flags().StringVarP(&listenAddr, "listen-address", "l", ":8080", "specify the address to listen for coordination instructions")
	return cmd
}

// func simulator(opts *simulatorOpts) {
// 	r := rand.New(rand.NewSource(0))

// 	mq := make(chan events.Event)
// 	go func() {
// 		for {
// 			msg := <-mq
// 			opts.appState.Logger.Debugf("queuing msg with id: '%s'", msg.GetMessageID())
// 			opts.Qimpl.Push(msg)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			time.Sleep(opts.Interval)
// 			opts.appState.Logger.Debug("generating event")
// 			msg := &events.EventMsg{
// 				ID: uuid.New().String(),
// 				Headers: map[string][]string{
// 					"User-Agent": {"test"},
// 					"Host":       {"test.com"},
// 				},
// 				URI:     "/test/1/2/3",
// 				Content: []byte(fmt.Sprintf("%d", r.Intn(100000))),
// 			}
// 			mq <- msg
// 		}
// 	}()
// }

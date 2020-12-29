package memory

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
)

// Relay is an in memory, reference implementation of the relay.
type Relay struct {
	relay.R
}

// Do will pass messages through an intermediary that may perform operations on the data.
func (r Relay) Do(ctx context.Context, outbound <-chan events.Event, inbound chan<- events.Event) error {
	// Pass all messages from the outbound queue to the inbound queue.
	for ctx.Err() == nil {
		select {
		case e := <-outbound:
			inbound <- e
		case _ = <-ctx.Done():
			break
		}
	}
	return ctx.Err()
}

package memory

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
	log "github.com/sirupsen/logrus"
)

// Relay is an in memory, reference implementation of the relay.
type Relay struct {
	relay.R
}

// ID returns a human readable identifier for this thing.
func (r Relay) ID() string {
	return "relay/memory"
}

// RDo will pass messages through an intermediary that may perform operations on the data.
func (r Relay) RDo(ctx context.Context, outbound <-chan []events.Event, inbound chan<- []events.Event) error {
	log.Info("starting in-memory relay")
	// Pass all messages from the outbound queue to the inbound queue.
	for ctx.Err() == nil {
		select {
		case e := <-outbound:
			log.Debugf("in-memory relay received '%d' events", len(e))
			inbound <- e
		case _ = <-ctx.Done():
			break
		}
	}
	return ctx.Err()
}

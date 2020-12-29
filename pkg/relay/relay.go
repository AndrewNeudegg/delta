package relay

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
)

// R (Relay) passes events from a source to a destination via any intermediary.
type R interface {
	// Do will pass messages through an intermediary that may perform operations on the data.
	// ctx: context will control the flow of the application.
	// outbound: will receive messages from the current applications.
	// inbound: will receive messages from the relay source.
	Do(ctx context.Context, outbound <-chan events.Event, inbound chan<- events.Event) error
}

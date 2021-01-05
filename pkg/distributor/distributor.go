package distributor

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
)

// D (distributor) will emit events to target services.
// It is responsible for identifying which services should be targetted
// and the method that those services should receive the events.
type D interface {
	DDo(context.Context, <-chan []events.Event) error // Do will emit events that are placed into the channel.
	ID() string                                       // ID returns a human readable identifier for this distributor.
}

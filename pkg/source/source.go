package source

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
)

// S is a generic source of events.
type S interface {
	SDo(context.Context, chan<- []events.Event) error // SDo will begin the loop for this source.
	ID() string                                       // ID returns a human readable identifier for this source.
}

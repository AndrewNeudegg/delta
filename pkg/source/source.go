package source

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
)

// S is a generic source of events.
type S interface {
	Do(context.Context, chan<- events.Event) error // Do will begin the loop for this source.
}

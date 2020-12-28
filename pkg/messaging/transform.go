package messaging

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
)

// Transform will apply modifications to an event.
type Transform interface {
	Do(context.Context, events.Event) (events.Event, error) // Do will apply the given transform.
}

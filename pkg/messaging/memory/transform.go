package memory

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/messaging"
)

// Transform is a noop transformer that will return the given event.
type Transform struct {
	messaging.Transform
}

// Do will return the given event with no modifications.
func (t Transform) Do(ctx context.Context, e events.Event) (events.Event, error) {
	return e, nil
}

package memory

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/messaging"
)

// Target is a specific in memory target for this event.
type Target struct {
	messaging.Target
	ch chan<- events.Event
}

// Send will emit an event towards the given target.
func (t Target) Send(ctx context.Context, e events.Event) error {
	t.ch <- e
	return nil
}

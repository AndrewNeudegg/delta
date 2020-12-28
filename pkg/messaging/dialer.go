package messaging

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
)

// Target defines how a msg can be routed to an application.
type Target interface {
	Send(context.Context, events.Event) error // Send will fire an event towards the selected target.
}

// Phonebook relates an event to a Target.
type Phonebook interface {
	Lookup(events.Event) (Target, error) // Lookup will return a target that is suitable for this event.
}

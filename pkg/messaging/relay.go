package messaging

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
)

// Relay is the relay interface, all relays should adhere to this.
type Relay interface {
	Do(context.Context, <-chan events.Event, Phonebook) error // Do will relay the given event to the next node.
}

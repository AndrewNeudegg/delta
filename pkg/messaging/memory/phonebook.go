package memory

import (
	"github.com/andrewneudegg/delta/pkg/messaging"
	"github.com/andrewneudegg/delta/pkg/events"
)

// Phonebook is a simple impl that will pass messages to the defined channel.
type Phonebook struct {
	messaging.Phonebook

	Ch chan<- events.Event
}

// Lookup will find a target suitable for the given event.
func (p Phonebook) Lookup(e events.Event) (Target, error) {
	return Target{
		ch: p.Ch,
	}, nil
}

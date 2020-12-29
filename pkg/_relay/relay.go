package _relay

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
)

// ------------ RELAY -------------

// R is the relay interface, all relays should adhere to this.
type R interface {
	Do(context.Context, <-chan events.Event)
}

// ------------ / RELAY -------------

// Server will serve events to requesting clients.
type Server interface {
	events.Distributor
	Start() error                   // Start the server.
	Stop(ctx context.Context) error // Stop the server.
}

// Client provides clients with a mechanism to access events.
type Client interface {
	GetNextEvent() (events.Event, error)        // GetNextEvent will return the next available event.
	GetNextNEvents(int) ([]events.Event, error) // GetNextNEvents will return up to N next events.
	CompleteEvent(events.Event) error           // CompleteEvent will complete the underlying event if required.
	CompleteNEvents([]events.Event) error       // CompleteNEvents will complete N underlying events if required.
	FailEvent(events.Event) error               // FailEvent will inform the RelayServer that the given event could not be processed.
	FailNEvents([]events.Event) error           // FailNEvents provides a convenience mechanism for bulk failing events.
	EstimateLoad() (int, error)                 // EstimateLoad provides an estimate for the amount of work being done by this server, between 0 and 100.
}

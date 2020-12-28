package sink

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
)

// ServerConfiguration exposes configuration applicable to all sinks
type ServerConfiguration struct {
	ToChan chan<- events.Event
}

// Server is the expectation for all sinks
type Server interface {
	Serve(ctx context.Context) error
	Stop(ctx context.Context) error
}

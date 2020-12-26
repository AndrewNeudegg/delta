package sink

import "context"

// ServerConfiguration exposes configuration applicable to all sinks
type ServerConfiguration struct {
	ToChan chan *SunkMessage
}

// Server is the expectation for all sinks
type Server interface {
	Serve(ctx context.Context) error
	Stop(ctx context.Context) error
}

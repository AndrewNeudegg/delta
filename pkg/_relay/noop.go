package _relay

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
)

// NoOpRelayOpts will write events to a channel.
type NoOpRelayOpts struct {
	Output chan<- events.Event
}

// noOpRelay is an in memory relay that just passes the incoming messages
// to a new channel that can be read by a distributor.
type noOpRelay struct {
	R
	config   *NoOpRelayOpts
	outputCh chan<- events.Event
}

func (n *noOpRelay) Do(ctx context.Context, msgCh <-chan events.Event) {
	for ctx.Err() == nil {
		msg := <-msgCh
		n.outputCh <- msg
	}
}

// NewNoOpRelay will return a relay that writes events back to a channel.
func NewNoOpRelay(c *NoOpRelayOpts) R {
	return &noOpRelay{
		config:   c,
		outputCh: c.Output,
	}
}

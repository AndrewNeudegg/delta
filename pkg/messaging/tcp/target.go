package tcp

import (
	"context"
	"encoding/gob"
	"net"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/messaging"
)

// Target is a specific in memory target for this event.
type Target struct {
	messaging.Target

	protocol string // tcp
	address  string // localhost:8080
}

// Send will emit an event towards the given target.
func (t Target) Send(ctx context.Context, e events.Event) error {
	conn, err := net.Dial(t.protocol, t.address)
	defer conn.Close()

	if err != nil {
		return err
	}

	// This is a little bit wasteful.
	encoder := gob.NewEncoder(conn)

	return encoder.Encode(e)
}

package memory

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/messaging"
)

// Relay passes event messages downstream.
type Relay struct {
	messaging.Relay
}

// Do will pass events from a channel to a target, by looking them up.
func (r Relay) Do(ctx context.Context, ch <-chan events.Event, p Phonebook) error {
	for ctx.Err() == nil {
		e := <-ch
		t, err := p.Lookup(e)
		if err != nil {
			return err
		}
		err = t.Send(ctx, e)
		if err != nil {
			return err
		}
	}

	return ctx.Err()
}

// Package stdout will distribute messages to this applications std out.
package stdout

import (
	"context"
	"fmt"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
)

// Distributor is a std out writer.
type Distributor struct {
	distributor.D
}

// Do will write all events to stdout.
func (d Distributor) Do(ctx context.Context, outbound <-chan events.Event) error {
	for ctx.Err() == nil {
		select {
		case e := <-outbound:
			fmt.Printf("Event: '%s'\n", e.GetMessageID())
		case _ = <-ctx.Done():
			break
		}
	}
	return ctx.Err()
}

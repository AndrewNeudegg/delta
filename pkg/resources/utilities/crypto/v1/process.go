package crypto1

import (
	"context"
	"fmt"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
)

// Process is simple noop.
type Process struct {
}

// ID defines what this thing is.
func (p Process) ID() string {
	return ID
}

// Type defines what type of resource this is.
func (p Process) Type() definitions.ResourceType {
	return definitions.ProcessType
}

// DoProcess will perform its function, on each collection placed into the channel
// eventually passing a similar collection to the output.
func (p Process) DoProcess(ctx context.Context, ch1 <-chan events.Collection, ch2 chan<- events.Collection) error {
	return fmt.Errorf("'%s' does not support process mode at the moment", ID)
}

package generators1

import (
	"context"
	"fmt"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
)

// Output is simple noop.
type Output struct {
	Configuration
}

// ID defines what this thing is.
func (o Output) ID() string {
	return ID
}

// Type defines what type of resource this is.
func (o Output) Type() definitions.ResourceType {
	return definitions.OutputType
}

// DoOutput will perform its function on each collection placed into the channel.
func (o Output) DoOutput(ctx context.Context, ch <-chan events.Collection) error {
	return fmt.Errorf("resource '%s' does not support output mode", ID)
}

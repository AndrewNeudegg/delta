package noop

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
)

// Output is simple noop.
type Output struct {
}

// ID defines what this thing is.
func (o Output) ID() string {
	return "resources/examples/simple"
}

// Type defines what type of resource this is.
func (o Output) Type() definitions.ResourceType {
	return definitions.OutputType
}

// DoOutput will perform its function on each collection placed into the channel.
func (o Output) DoOutput(ctx context.Context, ch <-chan events.Collection) error {
	<-ctx.Done()
	return nil
}

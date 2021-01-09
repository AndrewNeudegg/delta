package noop

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
)

// Input is simple noop.
type Input struct {
}

// ID defines what this thing is.
func (i Input) ID() string {
	return ID
}

// Type defines what type of resource this is.
func (i Input) Type() definitions.ResourceType {
	return definitions.InputType
}

// DoInput will accept collections of events, passing them into the channel.
func (i Input) DoInput(ctx context.Context, ch chan<- events.Collection) error {
	<-ctx.Done()
	return nil
}

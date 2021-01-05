package definitions

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
)

// ResourceType explains what a resource is at runtime.
type ResourceType string

const (
	// InputType of resource.
	InputType ResourceType = "input"

	// ProcessType of resource.
	ProcessType ResourceType = "process"

	// OutputType of resource.
	OutputType ResourceType = "output"

	// MetaType of resource.
	MetaType ResourceType = "meta"
)

// Resource is a thing that does things.
type Resource interface {
	ID() string         // ID defines what this thing is.
	Type() ResourceType // Type defines what type of resource this is.
}

// Input resource accepts events through some boundary.
type Input interface {
	Resource

	// DoInput will accept collections of events, passing them into the channel.
	DoInput(context.Context, chan<- events.Collection) error
}

// Process resource transfers or manipulates events.
type Process interface {
	Resource

	// DoProcess will perform its function, on each collection placed into the channel
	// eventually passing a similar collection to the output.
	DoProcess(context.Context, <-chan events.Collection, chan<- events.Collection) error
}

// Output resource passes events across some boundary.
type Output interface {
	Resource

	// DoOutput will perform its function on each collection placed into the channel.
	DoOutput(context.Context, <-chan events.Collection) error
}

// Meta resource augments one of the other resource types, or itself.
type Meta interface {
	Resource

	I([]Input) (Input, error)     // I, accepts inputs returning an input representing the given collection.
	P([]Process) (Process, error) // P, accepts process returning a process representing the given collection.
	O([]Output) (Output, error)   // O, accepts outputs returning an output representing the given collection.
	M([]Meta) (Meta, error)       // M, accepts meta resources returning a meta input for the given collection.
}

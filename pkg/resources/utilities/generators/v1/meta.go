package generators1

import (
	"fmt"

	"github.com/andrewneudegg/delta/pkg/resources/definitions"
)

// Meta is simple noop.
type Meta struct {
	Configuration
}

// ID defines what this thing is.
func (m Meta) ID() string {
	return ID
}

// Type defines what type of resource this is.
func (m Meta) Type() definitions.ResourceType {
	return definitions.MetaType
}

// I accepts inputs returning an input representing the given collection.
func (m Meta) I(i []definitions.Input) (definitions.Input, error) {
	return Input{}, fmt.Errorf("resource '%s' does not support meta mode", ID)
}

// P accepts process returning an process representing the given collection.
func (m Meta) P(p []definitions.Process) (definitions.Process, error) {
	return Process{}, fmt.Errorf("resource '%s' does not support meta mode", ID)
}

// O accepts outputs returning an output representing the given collection.
func (m Meta) O(o []definitions.Output) (definitions.Output, error) {
	return Output{}, fmt.Errorf("resource '%s' does not support meta mode", ID)
}

// M accepts meta resources returning a meta input for the given collection.
func (m Meta) M(me []definitions.Meta) (definitions.Meta, error) {
	return Meta{}, fmt.Errorf("resource '%s' does not support meta mode", ID)
}

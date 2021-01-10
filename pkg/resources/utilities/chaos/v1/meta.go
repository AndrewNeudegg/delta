package noop

import (
	"fmt"

	"github.com/andrewneudegg/delta/pkg/resources/definitions"
)

// Meta is simple noop.
type Meta struct {
	FailChance float32 `mapstructure:"failChance"`
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

	if len(i) != 1 {
		return nil, fmt.Errorf("'%s' only supports one input resource, found '%d'", ID, len(i))
	}

	return Input{
		i: i[0],
		failChance: m.FailChance,
	}, nil
}

// P accepts process returning an process representing the given collection.
func (m Meta) P(p []definitions.Process) (definitions.Process, error) {
	if len(p) != 1 {
		return nil, fmt.Errorf("'%s' only supports one process resource, found '%d'", ID, len(p))
	}

	return Process{
		p: p[0],
		failChance: m.FailChance,
	}, nil
}

// O accepts outputs returning an output representing the given collection.
func (m Meta) O(o []definitions.Output) (definitions.Output, error) {
	if len(o) != 1 {
		return nil, fmt.Errorf("'%s' only supports one output resource, found '%d'", ID, len(o))
	}

	return Output{
		o: o[0],
		failChance: m.FailChance,
	}, nil
}

// M accepts meta resources returning a meta input for the given collection.
func (m Meta) M(me []definitions.Meta) (definitions.Meta, error) {
	return Meta{}, fmt.Errorf("'%s' does not support meta mode", ID)
}

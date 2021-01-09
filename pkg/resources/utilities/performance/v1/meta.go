package performance1

import (
	"fmt"
	"time"

	"github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// Meta is simple noop.
type Meta struct {
	SampleWindow string `mapstructure:"sampleWindow"` // SampleWindow is the time that metrics are collected before reporting (i.e. 60s)
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
		return nil, fmt.Errorf("'%s' can only accept one input, got '%d'", ID, len(i))
	}

	dur, err := time.ParseDuration(m.SampleWindow)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse '%s' as duration", m.SampleWindow)
	}

	log.Debugf("returning meta input for '%s'", i[0].ID())

	return Input{
		i:            i[0],
		sampleWindow: dur,
	}, nil
}

// P accepts process returning an process representing the given collection.
func (m Meta) P(p []definitions.Process) (definitions.Process, error) {
	if len(p) != 1 {
		return nil, fmt.Errorf("'%s' can only accept one process, got '%d'", ID, len(p))
	}

	dur, err := time.ParseDuration(m.SampleWindow)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse '%s' as duration", m.SampleWindow)
	}

	return Process{
		p:            p[0],
		sampleWindow: dur,
	}, nil
}

// O accepts outputs returning an output representing the given collection.
func (m Meta) O(o []definitions.Output) (definitions.Output, error) {
	if len(o) != 1 {
		return nil, fmt.Errorf("'%s' can only accept one output, got '%d'", ID, len(o))
	}

	dur, err := time.ParseDuration(m.SampleWindow)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse '%s' as duration", m.SampleWindow)
	}

	return Output{
		o:            o[0],
		sampleWindow: dur,
	}, nil
}

// M accepts meta resources returning a meta input for the given collection.
func (m Meta) M(me []definitions.Meta) (definitions.Meta, error) {
	return Meta{}, fmt.Errorf("'%s' is not designed to return a meta resource", ID)
}

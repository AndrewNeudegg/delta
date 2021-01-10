package noop

import (
	"context"
	"fmt"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"

	log "github.com/sirupsen/logrus"
)

// Process is simple noop.
type Process struct {
	p definitions.Process

	failChance float32
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
	if p.p == nil {
		return fmt.Errorf("'%s' does not support input resource", ID)
	}
	log.Infof("starting '%s' DoProcess proxy for '%s'", ID, p.p.ID())


	return p.p.DoProcess(ctx, ch1, ch2)
}

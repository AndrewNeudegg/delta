package performance1

import (
	"context"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// Process is simple noop.
type Process struct {
	p            definitions.Process
	sampleWindow time.Duration
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
		return errors.Errorf("'%s' cannot be used as an process resource directly", ID)
	}
	log.Infof("starting '%s' DoProcess proxy for '%s'", ID, p.p.ID())

	// TODO!!

	<-ctx.Done()
	return nil
}

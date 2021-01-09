package performance1

import (
	"context"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/pkg/errors"
)

// Output is simple noop.
type Output struct {
	o            definitions.Output
	sampleWindow time.Duration
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
	if o.o == nil {
		return errors.Errorf("'%s' cannot be used as an output resource directly", ID)
	}

	<-ctx.Done()
	return nil
}

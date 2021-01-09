package console1

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"

	log "github.com/sirupsen/logrus"
)

// Output is simple noop.
type Output struct {
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
	log.Infof("'%s' awaiting events", ID)
	for {
		select {
		case eCol := <-ch:
			log.Infof("received '%d' events", len(eCol))
		case <-ctx.Done():
			return nil
		}
	}
}

package noop

import (
	"context"
	"fmt"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"

	log "github.com/sirupsen/logrus"
)

// Output is simple noop.
type Output struct {
	o definitions.Output

	failChance float32
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
		return fmt.Errorf("'%s' does not support input resource", ID)
	}
	log.Infof("starting '%s' DoOutput proxy for '%s'", ID, o.o.ID())

	proxyCh := make(chan events.Collection)

	go func(chIn chan events.Collection, chOut <-chan events.Collection) {
		for {
			select {
			case eCol := <-chOut:
				if isLucky(o.failChance) {
					chIn <- eCol
				} else {
					log.Debugf("'%s', '%d' events were unlucky and have been dropped", ID, len(eCol))
				}
			}
		}
	}(proxyCh, ch)

	return o.o.DoOutput(ctx, proxyCh)
}

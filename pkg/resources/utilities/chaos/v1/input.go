package noop

import (
	"context"
	"fmt"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"

	log "github.com/sirupsen/logrus"
)

// Input is simple noop.
type Input struct {
	i definitions.Input

	failChance float32
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
	if i.i == nil {
		return fmt.Errorf("'%s' does not support input resource", ID)
	}
	log.Infof("starting '%s' DoInput proxy for '%s'", ID, i.i.ID())

	proxyCh := make(chan events.Collection)

	go func(chIn chan events.Collection, chOut chan<- events.Collection) {
		for {
			select {
			case eCol := <-chIn:
				if isLucky(i.failChance) {
					chOut <- eCol
				}else{
					log.Debugf("'%s', '%d' events were unlucky and have been dropped", ID, len(eCol))
				}
			}
		}
	}(proxyCh, ch)

	return i.i.DoInput(ctx, proxyCh)
}

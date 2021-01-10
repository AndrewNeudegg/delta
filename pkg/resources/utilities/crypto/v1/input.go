package crypto1

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// Input is simple noop.
type Input struct {
	i definitions.Input

	c *crypto
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
		return errors.Errorf("'%s' cannot be used as an input resource directly", ID)
	}
	log.Infof("starting '%s' DoInput proxy for '%s'", ID, i.i.ID())

	proxyCh := make(chan events.Collection)

	go func(chIn chan events.Collection, chOut chan<- events.Collection) {
		for {
			select {
			case eCol := <-chIn:
				go func() {
					rCol, err := doCrypto(eCol, i.c)
					if err != nil {
						log.Error(err)
						return
					}
					chOut <- rCol
				}()
			}
		}
	}(proxyCh, ch)

	return i.i.DoInput(ctx, proxyCh)
}

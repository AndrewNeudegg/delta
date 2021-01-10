package crypto1

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// Output is simple noop.
type Output struct {
	o definitions.Output

	c *crypto
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
	log.Infof("starting '%s' DoOutput proxy for '%s'", ID, o.o.ID())

	proxyCh := make(chan events.Collection)

	go func(chIn <-chan events.Collection, chOut chan events.Collection) {
		for {
			select {
			case eCol := <-chIn:
				go func() {
					rCol, err := doCrypto(eCol, o.c)
					if err != nil {
						log.Error(err)
						return
					}
					chOut <- rCol
				}()
			}
		}
	}(ch, proxyCh)

	return o.o.DoOutput(ctx, proxyCh)
}

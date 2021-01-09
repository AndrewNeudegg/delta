package performance1

import (
	"context"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// Input is simple wrapper around another input resource.
type Input struct {
	i            definitions.Input
	sampleWindow time.Duration
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
	log.Infof("starting DoInput proxy for '%s'", i.i.ID())
	if i.i == nil {
		return errors.Errorf("'%s' cannot be used as an input resource directly", ID)
	}

	proxyCh := make(chan events.Collection)
	count := 0
	lastTime := time.Now()

	go func(chIn chan events.Collection, chOut chan<- events.Collection) {
		for {
			select {
			case eCol := <-chIn:
				count++
				chOut <- eCol
			}
		}
	}(proxyCh, ch)

	go func() {
		for {
			time.Sleep(i.sampleWindow)

			tDiff := time.Now().Sub(lastTime)
			metricFrame := float64(count) / tDiff.Seconds()
			log.Warnf("'%s' at '%f' tx/s (%d transactions / %f seconds)", i.i.ID(), metricFrame, count, tDiff.Seconds())

			lastTime = time.Now()
			count = 0
		}
	}()

	return i.i.DoInput(ctx, proxyCh)
}

package performance1

import (
	"context"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
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
	log.Infof("starting DoOutput proxy for '%s'", o.o.ID())

	if o.o == nil {
		return errors.Errorf("'%s' cannot be used as an output resource directly", ID)
	}

	proxyCh := make(chan events.Collection)
	count := 0
	lastTime := time.Now()

	go func(chIn <-chan events.Collection, chOut chan events.Collection) {
		for {
			select {
			case eCol := <-chIn:
				count++
				chOut <- eCol
			}
		}
	}(ch, proxyCh)

	go func() {
		for {
			time.Sleep(o.sampleWindow)

			tDiff := time.Now().Sub(lastTime)
			metricFrame := float64(count) / tDiff.Seconds()
			log.Warnf("'%s' at '%f' tx/s (%d transactions / %f seconds)", o.o.ID(), metricFrame, count, tDiff.Seconds())

			lastTime = time.Now()
			count = 0
		}
	}()

	return o.o.DoOutput(ctx, proxyCh)
}

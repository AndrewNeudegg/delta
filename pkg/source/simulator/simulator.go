package simulator

import (
	"context"
	"sync"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/source"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Source is a source that will generate events.
type Source struct {
	source.S

	Routines  int
	Interval  string
	Num       int
	BatchSize int
}

// ID is a human readable ID for this thing.
func (s Source) ID() string {
	return "source/simulator"
}

// SDo will start the simulator.
func (s Source) SDo(ctx context.Context, ch chan<- []events.Event) error {
	wg := sync.WaitGroup{}

	dur, err := time.ParseDuration(s.Interval)
	if err != nil {
		return errors.Wrapf(err, "could not parse '%s' as duration", s.Interval)
	}

	// retryF := func(ev events.Event, ch chan<- []events.Event) *func(error) {
	// 	f := func(e error) {
	// 		log.Debugf("event with id '%s' failed, retrying", ev.GetMessageID())
	// 		ch <- ev
	// 	}
	// 	return &f
	// }

	// completeF := func(e events.Event, ch chan<- []events.Event) *func() {
	// 	f := func() {
	// 		log.Debugf("event with id '%s' completed", e.GetMessageID())
	// 	}
	// 	return &f
	// }

	dF := func(dur time.Duration, ch chan<- []events.Event) {
		for ctx.Err() == nil {
			for i := 0; i < s.Num; i++ {
				eventCol := make([]events.Event, s.BatchSize)
				for i := 0; i < s.BatchSize; i++ {
					eventCol[i] = events.JunkEvent()
					
				}
				log.Debugf("generated '%d' event", len(eventCol))
				ch <- eventCol
			}
			time.Sleep(dur)
		}
		wg.Done()
	}

	for i := 0; i < s.Routines; i++ {
		wg.Add(1)
		go dF(dur, ch)
	}

	// do the serving.
	log.Infof("running event simulator with batch size '%d', '%d' goroutines and a delay of '%s'", s.Num, s.Routines, dur)
	wg.Wait()
	return ctx.Err()
}

package generators1

import (
	"context"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

const (
	// ID for this collection of resources.
	ID = "utilities/generators/v1"
)

// Configuration for this collection of resources.
type Configuration struct {
	Interval          string `mapstructure:"interval"`          // Interval specifies the delays between even
	NumberEvents      int    `mapstructure:"numberEvents"`      // NumberEvents is how many events will be in each `events.Collection`.
	NumberCollections int    `mapstructure:"numberCollections"` // NumberCollections is how many `events.Collection` will be delivered.
}

// RunGenerator actions the given config.
func RunGenerator(ctx context.Context, config Configuration, ch chan<- events.Collection) error {
	log.Infof("starting '%s'", ID)
	
	dur, err := time.ParseDuration(config.Interval)
	if err != nil {
		return errors.Wrapf(err, "could not parse '%s' as duration", config.Interval)
	}

	eCollections := make([]events.Collection, config.NumberCollections)
	for i := 0; i < config.NumberCollections; i++ {
		col := make(events.Collection, config.NumberEvents)
		for j := 0; j < config.NumberEvents; j++ {
			col[j] = events.JunkEvent()
		}
		eCollections[i] = col
	}

	for ctx.Err() == nil {
		for _, col := range eCollections {
			ch <- col
			log.Debugf("generated event collection with '%d' events", len(col))
		}
		time.Sleep(dur)
	}

	log.Infof("exiting generator with configuration: %+v", config)
	return ctx.Err()
}

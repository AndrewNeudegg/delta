package generators1

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"

	log "github.com/sirupsen/logrus"
)

// Input is simple event generator.
type Input struct {
	Interval          string `mapstructure:"interval"`          // Interval specifies the delays between even
	NumberEvents      int    `mapstructure:"numberEvents"`      // NumberEvents is how many events will be in each `events.Collection`.
	NumberCollections int    `mapstructure:"numberCollections"` // NumberCollections is how many `events.Collection` will be delivered.
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
	log.Infof("starting DoInput for '%s'", i.ID())

	return RunGenerator(ctx, Configuration{
		Interval:          i.Interval,
		NumberEvents:      i.NumberEvents,
		NumberCollections: i.NumberCollections,
	}, ch)
}

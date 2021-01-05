// Package performance will measure tx/s being received at the destination.
package performance

import (
	"context"
	"time"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Distributor is a std out writer.
type Distributor struct {
	distributor.D

	lastCheckTime       time.Time
	count               int
	MetricsPollInterval string `mapstructure:"metricsPollInterval"`
	metricFrames        []float64
}

// ID returns a human readable identifier for this thing.
func (d Distributor) ID() string {
	return "distributor/performance"
}

// DDo will write all events to stdout.
func (d *Distributor) DDo(ctx context.Context, outbound <-chan []events.Event) error {
	d.metricFrames = make([]float64, 0)
	d.lastCheckTime = time.Now()
	dur, err := time.ParseDuration(d.MetricsPollInterval)
	if err != nil {
		return errors.Wrapf(err, "could not parse '%s' as duration", d.MetricsPollInterval)
	}

	go func(dur time.Duration) {
		for ctx.Err() == nil {
			time.Sleep(dur)
			tDiff := time.Now().Sub(d.lastCheckTime)
			metricFrame := float64(d.count) / tDiff.Seconds()
			d.metricFrames = append(d.metricFrames, metricFrame)
			log.Warnf("Performing at '%f' tx/s (%d transactions / %f seconds)", metricFrame, d.count, tDiff.Seconds())
			d.count = 0
			d.lastCheckTime = time.Now()
		}
	}(dur)

	for {
		select {
		case eventCol := <-outbound:
			for _, e := range eventCol {
				d.count = d.count + 1
				e.Complete()
			}
		case _ = <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Package stdout will distribute messages to this applications std out.
package stdout

import (
	"context"
	"encoding/json"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Distributor is a std out writer.
type Distributor struct {
	distributor.D
}

// ID returns a human readable identifier for this thing.
func (d Distributor) ID() string {
	return "distributor/stdout"
}

// DDo will write all events to stdout.
func (d Distributor) DDo(ctx context.Context, outbound <-chan events.Event) error {
	for ctx.Err() == nil {
		select {
		case e := <-outbound:
			jsonBytes, err := json.Marshal(e)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to marshal json for distribution"))
			}
			log.Info(string(jsonBytes))
			e.Complete()
		case _ = <-ctx.Done():
			break
		}
	}
	return ctx.Err()
}

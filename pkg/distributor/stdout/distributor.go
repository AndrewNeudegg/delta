// Package stdout will distribute messages to this applications std out.
package stdout

import (
	"context"
	"encoding/json"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
	log "github.com/sirupsen/logrus"
)

// Distributor is a std out writer.
type Distributor struct {
	distributor.D
}

// Do will write all events to stdout.
func (d Distributor) Do(ctx context.Context, outbound <-chan events.Event) error {
	for ctx.Err() == nil {
		select {
		case e := <-outbound:
			jsonBytes, err := json.Marshal(e)
			if err != nil {
				log.Error(err)
			}
			log.Info(string(jsonBytes))
			e.Complete()
		case _ = <-ctx.Done():
			break
		}
	}
	return ctx.Err()
}

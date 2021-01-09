package http1

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"

	log "github.com/sirupsen/logrus"
)

// Output is simple noop.
type Output struct {
	TargetAddress string `mapstructure:"targetAddress"`

	client *http.Client
}

// ID defines what this thing is.
func (o Output) ID() string {
	return ID
}

// Type defines what type of resource this is.
func (o Output) Type() definitions.ResourceType {
	return definitions.OutputType
}

func (o Output) ship(eCol events.Collection) error {

	for _, e := range eCol {

		req, _ := http.NewRequest(
			"POST",
			fmt.Sprintf("%s%s", o.TargetAddress, e.GetURI()),
			bytes.NewBuffer(e.GetContent()))

		req.Header = e.GetHeaders()
		req.Header.Set("x-message-id", e.GetMessageID())
		req.Header.Set("Connection", "close")

		if _, err := o.client.Do(req); err != nil {
			return err
		}
		log.Debugf("'%s' emitted event with ID '%s'", ID, e.GetMessageID())
	}

	return nil
}

// DoOutput will perform its function on each collection placed into the channel.
func (o Output) DoOutput(ctx context.Context, ch <-chan events.Collection) error {


	o.client = &http.Client{
		Timeout: time.Second * 15,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}

	for {
		select {
		case eCol := <-ch:
			log.Debugf("'%s' received '%d' events for distribution", ID, len(eCol))
			go func ()  {
				err := o.ship(eCol)
				if err != nil {
					log.Error(err)
				}
			} ()
			
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

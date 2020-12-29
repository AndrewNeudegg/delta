package distributor

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
)

// Naive will pelt events at a single predefined address.
type Naive struct {
	Distributor
	Addr string // Addr to send events to (http://localhost:8080).
}

// Do will make a http post at the given Addr.
func (d Naive) Do(ctx context.Context, ch <-chan events.Event) error {

	// broadcast the event to the specific endpoint/s.
	broadcast := func(e events.Event) error {
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s%s", d.Addr, e.GetURI()),
			bytes.NewBuffer(e.GetContent()))

		if err != nil {
			return err
		}

		req.Header = e.GetHeaders()
		req.Header.Set("x-message-id", e.GetMessageID()) // a polite nod

		client := &http.Client{
			Timeout: time.Second * 15,
		}
		_, err = client.Do(req)
		return err
	}

	// backoffRetry will help when things get bumpy...
	backoffRetry := func(e events.Event) error {
		backoffSecs := time.Duration(1) * time.Second
		for i := 0; i < 5; i++ {
			time.Sleep(backoffSecs)

			if err := broadcast(e); err != nil {
				backoffSecs = backoffSecs * 2
				continue
			} else {
				return nil
			}

		}
		return fmt.Errorf("retry backoff reached")
	}

	for ctx.Err() == nil {
		// fan out immediately.
		go backoffRetry(<-ch)
	}

	return ctx.Err()
}

package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
)

// DirectDistributor will pelt events at a single predefined address.
type DirectDistributor struct {
	distributor.D
	Addr string // Addr to send events to (http://localhost:8080).
}

// ID returns a human readable identifier for this thing.
func (d DirectDistributor) ID() string {
	return "distributor/http/direct"
}

// DDo will make a http post at the given Addr.
func (d DirectDistributor) DDo(ctx context.Context, ch <-chan events.Event) error {

	client := &http.Client{
		Timeout: time.Second * 15,
	}

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

		_, err = client.Do(req)
		return err
	}

	// backoffRetry will help when things get bumpy...
	backoffRetry := func(e events.Event) error {
		backoffSecs := time.Duration(1) * time.Second
		for i := 0; i < 5; i++ {

			if err := broadcast(e); err != nil {
				backoffSecs = backoffSecs * 2
			} else {
				e.Complete()
				return nil
			}

			time.Sleep(backoffSecs)
		}
		err := fmt.Errorf("retry backoff reached")
		e.Fail(err)
		return err
	}

	for {
		// fan out immediately.
		select {
		case e := <-ch:
			go backoffRetry(e)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

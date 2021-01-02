package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			// DialContext: (&net.Dialer{
			// Timeout:   30 * time.Second,
			// KeepAlive: 30 * time.Second,
			// 	DualStack: true,
			// }).DialContext,
		},
	}

	// backoffRetry will help when things get bumpy...
	backoffRetry := func(e events.Event) error {
		backoffSecs := time.Duration(1) * time.Second
		for i := 0; i < 5; i++ {
			req, _ := http.NewRequest(
				"POST",
				fmt.Sprintf("%s%s", d.Addr, e.GetURI()),
				bytes.NewBuffer(e.GetContent()))
			req.Header = e.GetHeaders()
			req.Header.Set("x-message-id", e.GetMessageID())
			req.Header.Set("Connection", "close")

			if _, err := client.Do(req); err != nil {
				log.Error(errors.Wrap(err, "failed to do http request"))
				backoffSecs = backoffSecs * 2
			} else {
				log.Debugf("successfully sent event '%s'", e.GetMessageID())
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
		select {
		case e := <-ch:
			go backoffRetry(e)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

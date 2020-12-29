package naive

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/sink"
	"github.com/stretchr/testify/assert"
)

func TestNaiveSmoke(t *testing.T) {
	inboundEventsCh := make(chan events.Event)
	inboundEvents := []events.Event{}
	outboundEventsCh := make(chan events.Event)

	// Setup the dummy listener.
	server, _ := sink.NewHTTPSinkServer(&sink.HTTPSinkServerConfiguration{
		ServerConfiguration: sink.ServerConfiguration{
			ToChan: inboundEventsCh,
		},
		ListenAddr:  ":8057",
		MaxBodySize: 10000000,
	})
	go server.Serve(context.TODO())
	defer server.Stop(context.Background())
	go func(ch chan events.Event) {
		for {
			inboundEvents = append(inboundEvents, <-ch)
		}
	}(inboundEventsCh)

	// sleep for goroutines.
	time.Sleep(time.Second)

	// configure the distributor.
	n := Naive{
		Addr: "http://localhost:8057",
	}
	go n.Do(context.Background(), outboundEventsCh)

	// sleep for goroutines.
	time.Sleep(time.Second)

	numEvents := 10
	for i := 0; i < numEvents; i++ {
		count := fmt.Sprintf("%d", i)
		outboundEventsCh <- events.EventMsg{
			ID: count,
			Headers: map[string][]string{
				count: []string{count},
			},
			URI:     fmt.Sprintf("/%s", count),
			Content: []byte(count),
		}
	}

	// this test cpu starves the process somewhat
	// there is a race condition between emission and verification.
	// data ends up pooled in the channels so a sleep here allows
	// the test to iron out the inconsistencies before returning.
	time.Sleep(time.Second * 10)

	assert.Equal(t, numEvents, len(inboundEvents))
}

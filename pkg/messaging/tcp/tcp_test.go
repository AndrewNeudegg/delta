package tcp

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/stretchr/testify/assert"
)

func TestSmoke(t *testing.T) {

	// outCh is where messages that will be relayed go.
	outCh := make(chan events.Event)
	// inboundChan is where messages that have been received go.
	inboundChan := make(chan events.Event)
	inboundResults := []events.Event{}

	r := Relay{
		protocol:    "tcp",
		listenAddr:  ":8051",
		inboundChan: inboundChan,
	}

	go r.Do(context.Background(), outCh, Phonebook{
		KnownAddresses: []string{
			"localhost:8051",
		},
	})

	// start reading the arriving messages.
	go func() {
		for {
			e := <-inboundChan
			inboundResults = append(inboundResults, e)
		}
	}()

	// await goroutines
	time.Sleep(time.Second)

	numEvents := 10000
	for i := 0; i < numEvents; i++ {
		count := fmt.Sprintf("%d", i)
		outCh <- events.EventMsg{
			ID: count,
			Headers: map[string][]string{
				count: []string{count},
			},
			URI:     fmt.Sprintf("/%s", count),
			Content: []byte(count),
		}
	}

	// messages can be deep in the tcp stack, lets give them a moment to catchup.
	time.Sleep(time.Second)

	assert.Equal(t, numEvents, len(inboundResults))
}

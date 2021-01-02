package tcp

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/stretchr/testify/assert"
)

type router struct {
	protocol      string
	sendAddress   string
	listenAddress string
}

func (r router) Protocol() string {
	return r.protocol
}

func (r router) SendAddress() string {
	return r.sendAddress
}

func (r router) ListenAddress() string {
	return r.listenAddress
}

func TestSmoke(t *testing.T) {

	// outCh is where messages that will be relayed go.
	outboundCh := make(chan events.Event)
	// inboundChan is where messages that have been received go.
	inboundCh := make(chan events.Event)
	inboundResults := []events.Event{}

	r := Relay{
		Opts: router{
			protocol:      "tcp",
			sendAddress:   "localhost:8050",
			listenAddress: ":8050",
		},
	}

	go func() {
		err := r.RDo(context.Background(), outboundCh, inboundCh)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// start reading the arriving messages.
	go func() {
		for {
			e := <-inboundCh
			inboundResults = append(inboundResults, e)
		}
	}()

	numEvents := 100
	for i := 0; i < numEvents; i++ {
		count := fmt.Sprintf("%d", i)
		outboundCh <- events.EventMsg{
			ID: count,
			Headers: map[string][]string{
				count: []string{count},
			},
			URI:     fmt.Sprintf("/%s", count),
			Content: []byte(count),
		}
	}

	// messages can be deep in the tcp stack, lets give them a moment to catchup.
	time.Sleep(time.Second * 10)

	assert.Equal(t, numEvents, len(inboundResults))
}

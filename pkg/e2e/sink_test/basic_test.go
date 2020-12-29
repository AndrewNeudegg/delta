package sink_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay/memory"
	"github.com/andrewneudegg/delta/pkg/sink"
	"github.com/stretchr/testify/assert"
)

func TestE2EMemory(t *testing.T) {
	// Create a channel to pass messages from the api server.
	incomingEvents := make(chan events.Event)

	// Create a channel to collect 'distributed messages'.
	distributedMessageChan := make(chan events.Event)
	distributedMessages := []events.Event{}
	// pretend to be doing something with all these distributed events.
	go func(ch chan events.Event) {
		for {
			event := <-ch
			distributedMessages = append(distributedMessages, event)
		}
	}(distributedMessageChan)

	// relay will route messages from incomingEvents to distributedMessages.
	relay := memory.Relay{}
	go relay.Do(context.TODO(), incomingEvents, distributedMessageChan)

	// sinkServer will catch all events.
	sinkServer, _ := sink.NewHTTPSinkServer(&sink.HTTPSinkServerConfiguration{
		ServerConfiguration: sink.ServerConfiguration{
			ToChan: incomingEvents,
		},
		ListenAddr:  ":8090",
		MaxBodySize: 2097152,
	})
	go sinkServer.Serve(context.TODO())

	// pause for a bit in case we haven't context switches goroutines yet.
	time.Sleep(time.Second)

	// create a client.
	client := SinkClient{
		Addr: "http://localhost:8090",
	}

	// send a bunch of events.
	numEvents := 100000
	for i := 0; i < numEvents; i++ {
		result, err := client.Send("/test/hello", map[string][]string{
			"Host":         {fmt.Sprintf("%d.com", i)},
			"Content-Type": {"application/json"},
		}, []byte("hello"))
		assert.Nil(t, err)
		assert.NotEqual(t, "", result)
	}

	// check that we received the right number of events...
	assert.Equal(t, numEvents, len(distributedMessages))
}

// func TestSmoke(t *testing.T) {
// 	mq := make(chan events.Event)
// 	re := make(chan events.Event)
// 	resultantEvents := []events.Event{}
// 	numEvents := 100000

// 	relay := relay.NewNoOpRelay(&relay.NoOpRelayOpts{
// 		Output: re,
// 	})
// 	go relay.Do(context.TODO(), mq)
// 	go func() {
// 		for {
// 			msg := <-re
// 			resultantEvents = append(resultantEvents, msg)
// 		}
// 	}()

// 	sinkServer, _ := sink.NewHTTPSinkServer(&sink.HTTPSinkServerConfiguration{
// 		ServerConfiguration: sink.ServerConfiguration{
// 			ToChan: mq,
// 		},
// 		ListenAddr:  ":8090",
// 		MaxBodySize: 2097152,
// 	})
// 	go sinkServer.Serve(context.TODO())
// 	time.Sleep(time.Second)

// 	client := SinkClient{
// 		Addr: "http://localhost:8090",
// 	}

// 	for i := 0; i < numEvents; i++ {
// 		result, err := client.Send("/test/hello", map[string][]string{
// 			"Host":         {fmt.Sprintf("%d.com", i)},
// 			"Content-Type": {"application/json"},
// 		}, []byte("hello"))
// 		assert.Nil(t, err)
// 		assert.NotEqual(t, "", result)
// 	}

// 	assert.Equal(t, "/test/hello", resultantEvents[0].GetURI())
// 	assert.Equal(t, numEvents, len(resultantEvents))
// 	sinkServer.Stop(context.Background())
// }

package sink_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/distributor/naive"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
	"github.com/andrewneudegg/delta/pkg/relay/memory"
	"github.com/andrewneudegg/delta/pkg/source"
	"github.com/andrewneudegg/delta/pkg/source/sink/http"
	"github.com/stretchr/testify/assert"
)

func getSinkServer(listenAddr string, maxBodySize int) source.S {
	server := http.Sink{
		MaxBodySize: maxBodySize,
		ListenAddr:  listenAddr,
	}
	return server
}

func getMemoryRelay() relay.R {
	return memory.Relay{}
}

func getDistributor(sendAddr string) distributor.D {
	return naive.Naive{
		Addr: sendAddr,
	}
}

func TestSinkMemoryRelayDistributor(t *testing.T) {
	// This is an example of the entire application flow.
	// Messages are received and then redistributed.
	// Finally being passed to the exit points.

	// -- Entrypoint Server --
	iAddr := ":8080"
	iCh := make(chan events.Event)
	iS := getSinkServer(iAddr, 10000000)
	go iS.Do(context.TODO(), iCh)
	// -- / Entrypoint Server --

	// -- Endpoint Server --
	eAddr := ":7080"
	eCh := make(chan events.Event)
	rE := []events.Event{}
	eS := getSinkServer(eAddr, 10000000)
	go eS.Do(context.TODO(), eCh)
	go func(ch <-chan events.Event) {
		for {
			e := <-eCh
			rE = append(rE, e)
		}
	}(eCh)
	// -- / Endpoint Server --

	// -- Memory Relay --
	mR := getMemoryRelay()
	mRCh := make(chan events.Event)
	go mR.Do(context.TODO(), iCh, mRCh)
	// -- / Memory Relay --

	// -- Distributor --
	d := getDistributor("http://localhost:7080")
	go d.Do(context.TODO(), mRCh)
	// -- / Distributor --

	// ---------------------------------------------------

	client := SinkClient{
		Addr: "http://localhost:8080",
	}

	// send a bunch of events.
	numEvents := 1000
	f := func(i int) {
		result, err := client.Send("/test/hello", map[string][]string{
			"Host":         {fmt.Sprintf("%d.com", i)},
			"Content-Type": {"application/json"},
		}, []byte("hello"))
		assert.Nil(t, err)
		assert.NotEqual(t, "", result)
	}

	for i := 0; i < numEvents; i++ {
		f(i)
		time.Sleep(time.Microsecond)
	}

	time.Sleep(time.Second * 2)

	assert.Equal(t, numEvents, len(rE))
}

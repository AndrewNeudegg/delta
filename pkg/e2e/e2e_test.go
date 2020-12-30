package e2e

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"

	"testing"

	"github.com/andrewneudegg/delta/pkg/configuration"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/pipeline"
	"github.com/andrewneudegg/delta/pkg/source"
	"github.com/andrewneudegg/delta/pkg/source/sink/http"
)

func getSinkServer(listenAddr string, maxBodySize int) source.S {
	server := http.Sink{
		MaxBodySize: maxBodySize,
		ListenAddr:  listenAddr,
	}
	return server
}

func TestE2ESmoke(t *testing.T) {

	// endpoint server
	eAddr := ":5051"
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

	// ------

	config := []byte(`
applicationSettings: {}
sourceConfigurations:
  - name: sink/http
    config:
      ListenAddr: :5050
      MaxBodySize: 512
relayConfigs:
  - name: memory
    config: {}
distributorConfigurations:
  - name: naive
    config:
      Addr: http://localhost:5051
`)
	c := configuration.RawConfig{
		ConfigData: config,
	}
	configContainer, err := c.Load()
	_, err = pipeline.BuildPipeline(configContainer)
	assert.Nil(t, err)

	// ------
	client := SinkClient{
		Addr: "http://localhost:5050",
	}

	// send a bunch of events.
	numEvents := 2
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

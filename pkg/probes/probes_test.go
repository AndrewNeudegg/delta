package probes

import (
	"fmt"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_healthprobes(t *testing.T) {

	addr := ":6080"

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	probeServer := ProbeServer{
		ListenAddr: addr,
	}

	// server not launched yet
	resp, err := client.Get(fmt.Sprintf("http://localhost%s", addr))
	assert.NotNil(t, err, "error was nil but should have been timeout")

	// assert.Equal(t, 500, resp.StatusCode)
	go probeServer.StartProbeServer()
	runtime.Gosched()

	resp, err = client.Get(fmt.Sprintf("http://localhost%s/liveness", addr))
	assert.Nil(t, err, "error was not nil, but should have been")
	assert.Equal(t, 500, resp.StatusCode)

	resp, err = client.Get(fmt.Sprintf("http://localhost%s/readiness", addr))
	assert.Nil(t, err, "error was not nil, but should have been")
	assert.Equal(t, 500, resp.StatusCode)

	probeServer.AliveNow()

	resp, err = client.Get(fmt.Sprintf("http://localhost%s/liveness", addr))
	assert.Nil(t, err, "error was not nil, but should have been")
	assert.Equal(t, 200, resp.StatusCode)

	resp, err = client.Get(fmt.Sprintf("http://localhost%s/readiness", addr))
	assert.Nil(t, err, "error was not nil, but should have been")
	assert.Equal(t, 500, resp.StatusCode)

	probeServer.ReadyNow()

	resp, err = client.Get(fmt.Sprintf("http://localhost%s/liveness", addr))
	assert.Nil(t, err, "error was not nil, but should have been")
	assert.Equal(t, 200, resp.StatusCode)

	resp, err = client.Get(fmt.Sprintf("http://localhost%s/readiness", addr))
	assert.Nil(t, err, "error was not nil, but should have been")
	assert.Equal(t, 200, resp.StatusCode)
}

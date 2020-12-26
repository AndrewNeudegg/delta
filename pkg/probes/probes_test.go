package probes

import (
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_healthprobes(t *testing.T) {

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	probeServer := ProbeServer{
		ListenAddr: ":8080",
	}

	// server not launched yet
	resp, err := client.Get("http://localhost:8080")
	assert.NotNil(t, err, "error was nil but should have been timeout")

	// assert.Equal(t, 500, resp.StatusCode)
	go probeServer.StartProbeServer()
	runtime.Gosched()

	resp, err = client.Get("http://localhost:8080/liveness")
	assert.Nil(t, err, "error was not nil, but should have been")
	assert.Equal(t, 500, resp.StatusCode)

	resp, err = client.Get("http://localhost:8080/readiness")
	assert.Nil(t, err, "error was not nil, but should have been")
	assert.Equal(t, 500, resp.StatusCode)

	probeServer.AliveNow()

	resp, err = client.Get("http://localhost:8080/liveness")
	assert.Nil(t, err, "error was not nil, but should have been")
	assert.Equal(t, 200, resp.StatusCode)

	resp, err = client.Get("http://localhost:8080/readiness")
	assert.Nil(t, err, "error was not nil, but should have been")
	assert.Equal(t, 500, resp.StatusCode)

	probeServer.ReadyNow()

	resp, err = client.Get("http://localhost:8080/liveness")
	assert.Nil(t, err, "error was not nil, but should have been")
	assert.Equal(t, 200, resp.StatusCode)

	resp, err = client.Get("http://localhost:8080/readiness")
	assert.Nil(t, err, "error was not nil, but should have been")
	assert.Equal(t, 200, resp.StatusCode)
}

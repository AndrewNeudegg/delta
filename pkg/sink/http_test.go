package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

type dummyCounter struct {
	prometheus.Counter
}

func (d dummyCounter) Inc() {}

func getServer(mq chan<- events.Event, listenAddr string) Server {
	return &httpSinkServer{
		config: &HTTPSinkServerConfiguration{
			ServerConfiguration: ServerConfiguration{
				ToChan: mq,
			},
			ListenAddr:  listenAddr,
			MaxBodySize: 50, // 50 bytes
		},
		flowCounter:            dummyCounter{},
		flowCounterClientError: dummyCounter{},
		flowCounterServerError: dummyCounter{},
	}
}

func sendEvent(addr string, content []events.EventMsg) ([]events.Event, error) {
	resultantIDs := make([]events.Event, 0)
	client := &http.Client{}

	for _, v := range content {
		// override specifically for testing purposes
		targetAddr := fmt.Sprintf("http://localhost%s%s", addr, v.GetURI())

		req, err := http.NewRequest("POST", targetAddr, bytes.NewBuffer(v.GetContent()))
		if err != nil {
			return []events.Event{}, err
		}

		req.Header = v.GetHeaders()

		resp, err := client.Do(req)
		if err != nil {
			return []events.Event{}, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []events.Event{}, err
		}

		var result httpSinkServerResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			return []events.Event{}, err
		}

		v.ID = result.ID
		resultantIDs = append(resultantIDs, v)
	}

	return resultantIDs, nil
}

func TestSmoke(t *testing.T) {
	mq := make(chan events.Event)
	server := getServer(mq, ":8085")

	go server.Serve(context.TODO())
	defer server.Stop(context.TODO())
	time.Sleep(time.Second)
}

func TestSmokeFactory(t *testing.T) {
	mq := make(chan events.Event)
	server, err := NewHTTPSinkServer(&HTTPSinkServerConfiguration{
		ServerConfiguration: ServerConfiguration{
			ToChan: mq,
		},
		ListenAddr:  ":8085",
		MaxBodySize: 512,
	})
	assert.Nil(t, err)

	go server.Serve(context.TODO())
	defer server.Stop(context.TODO())
	time.Sleep(time.Second)
}

func TestSendEvent(t *testing.T) {
	addr := ":8085"
	inputData := []events.EventMsg{
		events.EventMsg{
			ID: "example",
			Headers: map[string][]string{
				"Content-Type": []string{"application/json"},
				"Host":         []string{"example.com"},
				"User-Agent":   []string{"example"},
			},
			URI:     "/test/hello",
			Content: []byte("hello world!"),
		},
	}

	mq := make(chan events.Event)

	sendResults := make([]events.Event, 0)
	go func() {
		for {
			msg := <-mq
			sendResults = append(sendResults, msg)
		}
	}()

	server := getServer(mq, addr)
	defer server.Stop(context.TODO())
	go server.Serve(context.TODO())
	time.Sleep(time.Second)

	sendResult, sendEventErr := sendEvent(addr, inputData)
	assert.Nil(t, sendEventErr, "did not expect an error sending event")
	assert.Equal(t, len(inputData), len(sendResult))
}

func TestSendEventWith2(t *testing.T) {
	addr := ":8085"
	inputData := []events.EventMsg{
		events.EventMsg{
			ID: "example",
			Headers: map[string][]string{
				"Content-Type": []string{"application/json"},
				"Host":         []string{"example.com"},
				"User-Agent":   []string{"example"},
			},
			URI:     "/test/hello1",
			Content: []byte("hello world!"),
		},
		events.EventMsg{
			ID: "example",
			Headers: map[string][]string{
				"Content-Type": []string{"application/json"},
				"Host":         []string{"example.com"},
				"User-Agent":   []string{"example"},
			},
			URI:     "/test/hello2",
			Content: []byte("hello world!"),
		},
	}

	mq := make(chan events.Event)

	sendResults := make([]events.Event, 0)
	go func() {
		for {
			msg := <-mq
			sendResults = append(sendResults, msg)
		}
	}()

	server := getServer(mq, addr)
	defer server.Stop(context.TODO())
	go server.Serve(context.TODO())
	time.Sleep(time.Second)

	sendResult, sendEventErr := sendEvent(addr, inputData)
	assert.Nil(t, sendEventErr, "did not expect an error sending event")
	assert.Equal(t, len(inputData), len(sendResult))

	server.Stop(context.TODO())
}

func TestGetOnDisallowedRoute(t *testing.T) {
	addr := ":8085"
	mq := make(chan events.Event)

	sendResults := make([]events.Event, 0)
	go func() {
		for {
			msg := <-mq
			sendResults = append(sendResults, msg)
		}
	}()

	server := getServer(mq, addr)
	defer server.Stop(context.TODO())
	go server.Serve(context.TODO())

	time.Sleep(time.Second)
	resp, err := http.Get("http://localhost:8085/hello/testing")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestVeryLargeBody(t *testing.T) {
	addr := ":8085"
	inputData := []events.EventMsg{
		events.EventMsg{
			ID: "example",
			Headers: map[string][]string{
				"Content-Type": []string{"application/json"},
				"Host":         []string{"example.com"},
				"User-Agent":   []string{"example"},
			},
			URI:     "/test/hello1",
			Content: []byte("hello world!hello world!hello world!hello world!hello world!hello world!hello world!hello world!hello world!"),
		},
	}

	mq := make(chan events.Event)

	sendResults := make([]events.Event, 0)
	go func() {
		for {
			msg := <-mq
			sendResults = append(sendResults, msg)
		}
	}()

	server := getServer(mq, addr)
	defer server.Stop(context.TODO())
	go server.Serve(context.TODO())
	time.Sleep(time.Second)

	sendResult, sendEventErr := sendEvent(addr, inputData)
	assert.NotNil(t, sendEventErr, "expected an error sending event")
	assert.Equal(t, 0, len(sendResult))
}

func str2ptr(a string) *string {
	return &a
}

func str2ptrByte(a string) *[]byte {
	thing := []byte(a)
	return &thing
}

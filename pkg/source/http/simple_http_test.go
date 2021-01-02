package http

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
	"github.com/andrewneudegg/delta/pkg/source"
	"github.com/stretchr/testify/assert"
)

// newServer will get a new sink server.
func newServer(listenAddr string, maxBodySize int) source.S {
	return &SimpleHTTPSink{
		ListenAddr:  listenAddr,
		MaxBodySize: maxBodySize,
	}
}

func sendEvent(addr string, content []events.EventMsg) ([]events.Event, error) {
	resultantIDs := make([]events.Event, 0)
	client := &http.Client{
		Timeout: time.Second * 15,
	}

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
	ch := make(chan events.Event)
	server := newServer(":8074", 512)

	go server.SDo(context.TODO(), ch)
	time.Sleep(time.Second)
}

func TestSendEvent(t *testing.T) {
	addr := ":8073"

	// ------- Setup -------
	ch := make(chan events.Event)
	sendResults := make([]events.Event, 0)
	go func(ch chan events.Event) {
		for {
			msg := <-ch
			msg.Complete()
			sendResults = append(sendResults, msg)
		}
	}(ch)
	time.Sleep(time.Second)

	server := newServer(addr, 512)
	go server.SDo(context.TODO(), ch)
	time.Sleep(time.Second)
	// ------- Test --------

	inputData := []events.EventMsg{
		{
			ID: "example",
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
				"Host":         {"example.com"},
				"User-Agent":   {"example"},
			},
			URI:     "/test/hello",
			Content: []byte("hello world!"),
		},
	}

	sendResult, sendEventErr := sendEvent(addr, inputData)
	if sendEventErr != nil {
		t.Fatal(sendEventErr)
	}

	time.Sleep(time.Second * 1)
	assert.Equal(t, len(inputData), len(sendResult))
}

func TestSendEventWith2(t *testing.T) {
	addr := ":8072"

	ch := make(chan events.Event)
	server := newServer(addr, 512)
	go func() {
		err := server.SDo(context.TODO(), ch)
		assert.Nil(t, err)
	}()

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

	sendResults := make([]events.Event, 0)
	go func() {
		for {
			msg := <-ch
			msg.Complete() // must do this..
			sendResults = append(sendResults, msg)
		}
	}()

	time.Sleep(time.Second)

	sendResult, sendEventErr := sendEvent(addr, inputData)
	assert.Nil(t, sendEventErr, "did not expect an error sending event")
	assert.Equal(t, len(inputData), len(sendResult))
}

func TestGetOnDisallowedRoute(t *testing.T) {
	addr := ":8075"

	ch := make(chan events.Event)
	server := newServer(addr, 512)
	go func() {
		err := server.SDo(context.TODO(), ch)
		assert.Nil(t, err)
	}()

	sendResults := make([]events.Event, 0)
	go func() {
		for {
			msg := <-ch
			msg.Complete()
			sendResults = append(sendResults, msg)
		}
	}()

	time.Sleep(time.Second)
	resp, err := http.Get("http://localhost:8075/hello/testing")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestVeryLargeBody(t *testing.T) {
	addr := ":8076"

	ch := make(chan events.Event)
	server := newServer(addr, 10)
	go func() {
		err := server.SDo(context.TODO(), ch)
		assert.Nil(t, err)
	}()

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

	sendResults := make([]events.Event, 0)
	go func() {
		for {
			msg := <-ch
			msg.Complete()
			sendResults = append(sendResults, msg)
		}
	}()

	time.Sleep(time.Second)

	sendResult, sendEventErr := sendEvent(addr, inputData)
	assert.NotNil(t, sendEventErr, "expected an error sending event")
	assert.Equal(t, 0, len(sendResult))
}

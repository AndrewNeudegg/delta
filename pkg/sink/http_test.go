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

	"github.com/stretchr/testify/assert"
)

func getServer(mq chan<- *SunkMessage, listenAddr string) Server {
	return &httpSinkServer{
		config: &HTTPSinkServerConfiguration{
			ServerConfiguration: ServerConfiguration{
				ToChan: mq,
			},
			ListenAddr: listenAddr,
		},
	}
}

func sendEvent(addr string, content []SunkMessage) ([]SunkMessage, error) {
	resultantIDs := make([]SunkMessage, 0)
	client := &http.Client{}

	for _, v := range content {
		// override specifically for testing purposes
		targetAddr := fmt.Sprintf("http://localhost%s%s", addr, *v.URI)

		req, err := http.NewRequest("POST", targetAddr, bytes.NewBuffer(*v.Content))
		if err != nil {
			return []SunkMessage{}, err
		}

		req.Header.Set("Content-Type", *v.ContentType)
		req.Header.Set("User-Agent", *v.UserAgent)
		req.Header.Set("Host", *v.Host)

		resp, err := client.Do(req)
		if err != nil {
			return []SunkMessage{}, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []SunkMessage{}, err
		}

		var result httpSinkServerResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			return []SunkMessage{}, err
		}

		v.MessageID = &result.ID
		resultantIDs = append(resultantIDs, v)
	}

	return resultantIDs, nil
}

func TestSmoke(t *testing.T) {
	mq := make(chan *SunkMessage)
	server := getServer(mq, ":8085")

	go server.Serve(context.TODO())
	defer server.Stop(context.TODO())
	time.Sleep(time.Second)
}

func TestSendEvent(t *testing.T) {
	addr := ":8085"
	inputData := []SunkMessage{
		{
			MessageID:   str2ptr("empty"),
			Host:        str2ptr("example.com"),
			ContentType: str2ptr("application/json"),
			UserAgent:   str2ptr("testing"),
			URI:         str2ptr("/test/hello"),
			Content:     str2ptrByte("hello world"),
		},
	}

	mq := make(chan *SunkMessage)

	sendResults := make([]SunkMessage, 0)
	go func() {
		for {
			msg := <-mq
			sendResults = append(sendResults, *msg)
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

func TestSendEventWith2(t *testing.T) {
	addr := ":8085"
	inputData := []SunkMessage{
		{
			MessageID:   str2ptr("empty"),
			Host:        str2ptr("example.com"),
			ContentType: str2ptr("application/json"),
			UserAgent:   str2ptr("testing"),
			URI:         str2ptr("/test/hello"),
			Content:     str2ptrByte("hello world"),
		},
		{
			MessageID:   str2ptr("empty"),
			Host:        str2ptr("example.com"),
			ContentType: str2ptr("application/json"),
			UserAgent:   str2ptr("testing"),
			URI:         str2ptr("/test/hello"),
			Content:     str2ptrByte("hello world"),
		},
	}

	mq := make(chan *SunkMessage)

	sendResults := make([]SunkMessage, 0)
	go func() {
		for {
			msg := <-mq
			sendResults = append(sendResults, *msg)
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
	mq := make(chan *SunkMessage)

	sendResults := make([]SunkMessage, 0)
	go func() {
		for {
			msg := <-mq
			sendResults = append(sendResults, *msg)
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

func str2ptr(a string) *string {
	return &a
}

func str2ptrByte(a string) *[]byte {
	thing := []byte(a)
	return &thing
}

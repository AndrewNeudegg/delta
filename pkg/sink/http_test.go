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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

type dummyCounter struct {
	prometheus.Counter
}

func (d dummyCounter) Inc() {}

func getServer(mq chan<- *SunkMessage, listenAddr string) Server {
	return &httpSinkServer{
		config: &HTTPSinkServerConfiguration{
			ServerConfiguration: ServerConfiguration{
				ToChan: mq,
			},
			ListenAddr:  listenAddr,
			MaxBodySize: 512, // 512 bytes
		},
		flowCounter: dummyCounter{},
		flowCounterClientError: dummyCounter{},
		flowCounterServerError: dummyCounter{},
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

func TestSmokeFactory(t *testing.T) {
	mq := make(chan *SunkMessage)
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

func TestVeryLargeBody(t *testing.T) {
	addr := ":8085"
	inputData := []SunkMessage{
		{
			MessageID:   str2ptr("empty"),
			Host:        str2ptr("example.com"),
			ContentType: str2ptr("application/json"),
			UserAgent:   str2ptr("testing"),
			URI:         str2ptr("/test/hello"),
			Content:     str2ptrByte("70sAESz3wsZnqIp4tZ6sImidbjXbjBbYFmcayJDZC0GgyViA51jrWDIM0ePkS5RoA9SqhmoPkIoFy6cPw2DmINc7dby1gsXdWgZ33JoMzecz3Mmk2UDsLulfmrlEuYa7IEXLB34fx7pkCsm9NxjP1v6sRSp9IXSjw8W4Jo4Cc2KeIYkECW3YP71Za7YznGXUHyeueP6qJ3MHgDUWutqRVhuG6wj7xR8rTFbVrFB7GLsqtuVQ7j6f4dkmOvDueh0EYA0uAq5we3hxI4eE1EguXe7y0EPr9FO93UcjzAIrT5thHBLnrQFBBHSzkpx04h84yRKjMhrfEN5JkxKe7MZzCiazNUcivmuGbrsh2aTsZEUVYBH8qZMxn3pBDQJK9vE38kV6Ew4Yyv6Z2eqC25ViguQ6PsNLzjo91mvJFnXSnZbPCiabXf68Vz3DBjVNxFM4uQuIkjG3le3To3EZiHKCDJ4uinD0ftyd31LBPfwiHO8SIfZMoxUqynRLhOaOXZG7LaM"),
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

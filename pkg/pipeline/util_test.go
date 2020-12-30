package pipeline

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/stretchr/testify/assert"
)

func TestFanInSmoke(t *testing.T) {

	resultCh := make(chan events.Event)
	results := []events.Event{}

	go func() {
		for {
			e := <-resultCh
			results = append(results, e)
		}
	}()

	fans := []chan events.Event{}
	numFans := 10
	for i := 0; i < numFans; i++ {
		fans = append(fans, make(chan events.Event))
	}

	go fanIn(context.TODO(), fans, resultCh)

	for i := 0; i < numFans; i++ {
		count := fmt.Sprintf("%d", i)

		fans[i] <- events.EventMsg{
			ID: count,
			Headers: map[string][]string{
				count: []string{count},
			},
			URI:     fmt.Sprintf("/%s", count),
			Content: []byte(count),
		}
	}

	time.Sleep(time.Second)

	assert.Equal(t, numFans, len(results))
}

func TestFanOutSmoke(t *testing.T) {

	inputCh := make(chan events.Event)

	// -- Test Setup

	resultCh := make(chan events.Event)
	results := []events.Event{}

	go func() {
		for {
			e := <-resultCh
			results = append(results, e)
		}
	}()

	fans := []chan events.Event{}
	numFans := 10
	for i := 0; i < numFans; i++ {
		fans = append(fans, make(chan events.Event))
	}

	go fanIn(context.TODO(), fans, resultCh)

	// -- Test --

	go fanOut(context.TODO(), inputCh, fans)

	for i := 0; i < numFans; i++ {
		count := fmt.Sprintf("%d", i)

		inputCh <- events.EventMsg{
			ID: count,
			Headers: map[string][]string{
				count: []string{count},
			},
			URI:     fmt.Sprintf("/%s", count),
			Content: []byte(count),
		}
	}

	time.Sleep(time.Second)

	assert.Equal(t, numFans*numFans, len(results))
}

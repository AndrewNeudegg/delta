package utils

import (
	"context"
	"sync"

	"github.com/andrewneudegg/delta/pkg/events"
)

// Channels is a generic container struct.
type Channels struct{}

// FanIn will merge multiple input channels to a singular output channel.
func (c *Channels) FanIn(ctx context.Context, chs []chan []events.Event, combined chan []events.Event) error {
	wg := sync.WaitGroup{}

	merge := func(ch <-chan []events.Event) {
		for {
			select {
			case e := <-ch:
				combined <- e
			case _ = <-ctx.Done():
				wg.Done()
				return
			}
		}
	}

	for _, v := range chs {
		wg.Add(1)
		go merge(v)
	}

	wg.Wait()
	return ctx.Err()
}

// FanOut will split a single input channel into multiple output channels.
func (c *Channels) FanOut(ctx context.Context, ch chan []events.Event, outputs []chan []events.Event) error {
	dispatch := func(chs []chan []events.Event, e []events.Event) {
		for _, v := range chs {
			go func(to chan []events.Event, e []events.Event) {
				to <- e
			}(v, e)
		}
	}

	for {
		select {
		case e := <-ch:
			go dispatch(outputs, e)
		case _ = <-ctx.Done():
			return ctx.Err()
		}
	}
}

package utils

import (
	"context"
	"sync"

	"github.com/andrewneudegg/delta/pkg/events"
)

type Channels struct{}

// FanIn will merge multiple input channels to a singular output channel.
func (c *Channels) FanIn(ctx context.Context, chs []chan events.Event, combined chan events.Event) error {
	wg := sync.WaitGroup{}

	merge := func(ch <-chan events.Event) {
		for ctx.Err() == nil {
			select {
			case e := <-ch:
				combined <- e
			case _ = <-ctx.Done():
				break
			}
		}
		wg.Done()
	}

	for _, v := range chs {
		wg.Add(1)
		go merge(v)
	}

	wg.Wait()
	return ctx.Err()
}

// FanOut will split a single input channel into multiple output channels.
func (c *Channels) FanOut(ctx context.Context, ch chan events.Event, outputs []chan events.Event) error {

	write := func(ch chan events.Event, e events.Event) {
		ch <- e
	}

	dispatch := func(chs []chan events.Event, e events.Event) {
		for _, v := range chs {
			go write(v, e)
		}
	}

	for ctx.Err() == nil {
		select {
		case e := <-ch:
			go dispatch(outputs, e)
		case _ = <-ctx.Done():
			break
		}
	}

	return ctx.Err()
}

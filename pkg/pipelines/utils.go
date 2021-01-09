package pipelines

import (
	"context"
	"sync"

	"github.com/andrewneudegg/delta/pkg/events"
)

// FanIn will merge multiple input channels to a singular output channel.
func FanIn(ctx context.Context, chs []chan events.Collection, combined chan events.Collection) error {
	wg := sync.WaitGroup{}

	merge := func(ch <-chan events.Collection) {
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
func FanOut(ctx context.Context, ch chan events.Collection, outputs []chan events.Collection) error {
	dispatch := func(chs []chan events.Collection, e events.Collection) {
		for _, v := range chs {
			go func(to chan events.Collection, e events.Collection) {
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

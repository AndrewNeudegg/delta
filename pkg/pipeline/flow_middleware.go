package pipeline

import "github.com/andrewneudegg/delta/pkg/events"

// EventMiddlewareFunc is a function that will be called on each event.
type EventMiddlewareFunc func([]events.Event)

// Inject wraps the given channel with another channel that actions the middleware on each item in a separate goroutine.
func Inject(ch chan []events.Event, middleware EventMiddlewareFunc) chan []events.Event {
	newOutPipe := make(chan []events.Event)

	go func() {
		for {
			select {
			case e := <-ch:
				go middleware(e)
				newOutPipe <- e
			}
		}
	}()

	return newOutPipe
}

// // NoopEventMiddleware does not do anything, but matches the prometheus definition.
// func NoopEventMiddleware(counterName string) func(events.Event) {
// 	return func(e events.Event) {}
// }

// // PrometheusTelemetryEventMiddleware will increment a counter with the given name.
// func PrometheusTelemetryEventMiddleware(counterName string) func(events.Event) {

// 	counter := promauto.NewCounter(prometheus.CounterOpts{
// 		Name: counterName,
// 		Help: fmt.Sprintf("autogenerated counter for '%s'", counterName),
// 	})

// 	return func(e events.Event) {
// 		counter.Add(1)
// 	}
// }

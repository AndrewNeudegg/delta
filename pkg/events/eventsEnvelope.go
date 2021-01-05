package events

import "github.com/pkg/errors"

// Collection of events.
type Collection []Event

// Complete these events.
func (c Collection) Complete() {
	for _, e := range c {
		e.Complete()
	}
}

// Fail these events.
func (c Collection) Fail(err error) {
	for _, e := range c {
		e.Fail(err)
	}
}

// Await all events or error.
func (c Collection) Await() error {
	var aErr error
	for _, e := range c {
		err := e.Await()
		if err != nil {
			if aErr == nil {
				aErr = err
			}
			aErr = errors.Wrapf(aErr, "event '%s' could not be awaited", e.GetMessageID())
		}
	}
	return aErr
}

package example

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/meta"
	"github.com/andrewneudegg/delta/pkg/relay"
	"github.com/andrewneudegg/delta/pkg/source"
)

// Example doesn't really do anything here.
type Example struct {
	meta.M
}

// DoS will do S with some modification.
func (m Example) DoS(ctx context.Context, ch chan events.Event, s source.S) error {
	return s.Do(ctx, ch)
}

// DoR will do R with some modification.
func (m Example) DoR(ctx context.Context, chOut chan events.Event, chIn chan events.Event, r relay.R) error {
	return r.Do(ctx, chOut, chIn)
}

// DoD will do D with some modification.
func (m Example) DoD(ctx context.Context, ch chan events.Event, d distributor.D) error {
	return d.Do(ctx, ch)
}

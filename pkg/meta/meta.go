package meta

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
	"github.com/andrewneudegg/delta/pkg/source"
)

// M defines a block that can augment other functions.
type M interface {
	DoS(context.Context, chan events.Event, source.S) error
	DoR(context.Context, chan events.Event, chan events.Event, relay.R) error
	DoD(context.Context, chan events.Event, distributor.D) error
}

package meta

import (
	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/relay"
	"github.com/andrewneudegg/delta/pkg/source"
)

// M defines a block that can augment other functions.
type M interface {
	S([]source.S) (source.S, error)           // S, source intermediary.
	R([]relay.R) (relay.R, error)             // R, relay intermediary.
	D([]distributor.D) (distributor.D, error) // D, distributor intermediary.
	ID() string                               // ID returns a human readable identifier for this Meta.
}

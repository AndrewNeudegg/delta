package builder

import (
	"github.com/mitchellh/mapstructure"

	"github.com/andrewneudegg/delta/pkg/relay"
	"github.com/andrewneudegg/delta/pkg/relay/memory"
	"github.com/andrewneudegg/delta/pkg/relay/tcp"
)

// Get will return the given source with its data values initialised.
func Get(distributorName string, distributorConfiguration interface{}) (relay.R, error) {
	switch distributorName {
	case "memory":
		r := memory.Relay{}
		err := mapstructure.Decode(distributorConfiguration, &r)
		return r, err
	case "tcp":
		r := tcp.Relay{}
		err := mapstructure.Decode(distributorConfiguration, &r)
		return r, err
	}

	return nil, nil
}

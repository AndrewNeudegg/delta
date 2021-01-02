package builder

import (
	"github.com/andrewneudegg/delta/pkg/meta"
	"github.com/andrewneudegg/delta/pkg/meta/chaos"

	"github.com/mitchellh/mapstructure"
)

// Get will return the given source with its data values initialised.
func Get(distributorName string, metaConfiguration interface{}) (meta.M, error) {
	switch distributorName {
	case "meta/chaos/simple":
		m := chaos.Simple{}
		err := mapstructure.Decode(metaConfiguration, &m)
		return &m, err
	}

	return nil, nil
}

package builder

import (
	"github.com/andrewneudegg/delta/pkg/meta"
	"github.com/andrewneudegg/delta/pkg/meta/chaos"
	"github.com/andrewneudegg/delta/pkg/meta/example"

	"github.com/mitchellh/mapstructure"
)

// Get will return the given source with its data values initialised.
func Get(distributorName string, metaConfiguration interface{}) (meta.M, error) {
	switch distributorName {
	case "meta/example":
		m := example.Example{}
		err := mapstructure.Decode(metaConfiguration, &m)
		return m, err
	case "meta/example2":
		m := example.Example{}
		err := mapstructure.Decode(metaConfiguration, &m)
		return m, err
	case "meta/chaos/simple":
		m := chaos.ChaosSimple{}
		err := mapstructure.Decode(metaConfiguration, &m)
		return m, err
	}

	return nil, nil
}

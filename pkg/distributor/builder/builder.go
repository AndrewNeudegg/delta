package builder

import (
	"github.com/mitchellh/mapstructure"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/distributor/http"
	"github.com/andrewneudegg/delta/pkg/distributor/stdout"
)

// Get will return the given source with its data values initialised.
func Get(distributorName string, distributorConfiguration interface{}) (distributor.D, error) {
	switch distributorName {
	case "http/direct":
		d := http.DirectDistributor{}
		err := mapstructure.Decode(distributorConfiguration, &d)
		return d, err
	case "stdout":
		d := stdout.Distributor{}
		err := mapstructure.Decode(distributorConfiguration, &d)
		return d, err
	}

	return nil, nil
}

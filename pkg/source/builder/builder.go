package builder

import (
	"github.com/mitchellh/mapstructure"

	"github.com/andrewneudegg/delta/pkg/source"
	"github.com/andrewneudegg/delta/pkg/source/http"
	"github.com/andrewneudegg/delta/pkg/source/simulator"
)

// Get will return the given source with its data values initialised.
func Get(sourceName string, sourceConfiguration interface{}) (source.S, error) {
	switch sourceName {
	case "http/simple":
		source := http.SimpleHTTPSink{}
		err := mapstructure.Decode(sourceConfiguration, &source)
		return source, err
	case "source/simulator/simple":
		source := simulator.Source{}
		err := mapstructure.Decode(sourceConfiguration, &source)
		return source, err
	}

	return nil, nil
}

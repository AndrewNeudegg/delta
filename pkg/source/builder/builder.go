package builder

import (
	"github.com/mitchellh/mapstructure"

	"github.com/andrewneudegg/delta/pkg/source"
	"github.com/andrewneudegg/delta/pkg/source/http"
)

// Get will return the given source with its data values initialised.
func Get(sourceName string, sourceConfiguration interface{}) (source.S, error) {
	switch sourceName {
	case "http/simple":
		source := http.SimpleHttpSink{}
		err := mapstructure.Decode(sourceConfiguration, &source)
		return source, err
	}

	return nil, nil
}

package pipelines

import (
	"fmt"

	"github.com/andrewneudegg/delta/pkg/pipelines/definitions"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// BuildPipeline will construct a pipeline.
func BuildPipeline(identifier string, data interface{}, m map[string]definitions.Pipeline) (definitions.Pipeline, error) {
	val, ok := m[identifier]
	if !ok {
		return nil, fmt.Errorf("source '%s' is unknown", identifier)
	}

	if err := mapstructure.Decode(data, &val); err != nil {
		return nil, errors.Wrapf(err, "could not decode mapstructure for '%s'", identifier)
	}

	return val, nil
}

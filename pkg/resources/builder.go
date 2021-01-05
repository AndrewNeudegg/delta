package resources

import (
	"fmt"

	"github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/andrewneudegg/delta/pkg/resources/examples/noop"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// BuildInputResource will construct a resource from the given inputs.
func BuildInputResource(identifier string, data interface{}) (definitions.Input, error) {
	// Duplicate names are prohibited.
	m := map[string]definitions.Input{
		noop.Input{}.ID(): noop.Input{},
	}

	val, ok := m[identifier]
	if !ok {
		return nil, fmt.Errorf("source '%s' is unknown", identifier)
	}

	if err := mapstructure.Decode(data, &val); err != nil {
		return nil, errors.Wrapf(err, "could not decode mapstructure for '%s'", identifier)
	}

	return val, nil
}

// BuildProcessResource will construct a resource from the given process.
func BuildProcessResource(identifier string, data interface{}) (definitions.Process, error) {
	// Duplicate names are prohibited.
	m := map[string]definitions.Process{
		noop.Process{}.ID(): noop.Process{},
	}

	val, ok := m[identifier]
	if !ok {
		return nil, fmt.Errorf("source '%s' is unknown", identifier)
	}

	if err := mapstructure.Decode(data, &val); err != nil {
		return nil, errors.Wrapf(err, "could not decode mapstructure for '%s'", identifier)
	}

	return val, nil
}

// BuildOutputResource will construct a resource from the given outputs.
func BuildOutputResource(identifier string, data interface{}) (definitions.Output, error) {
	// Duplicate names are prohibited.
	m := map[string]definitions.Output{
		noop.Output{}.ID(): noop.Output{},
	}

	val, ok := m[identifier]
	if !ok {
		return nil, fmt.Errorf("source '%s' is unknown", identifier)
	}

	if err := mapstructure.Decode(data, &val); err != nil {
		return nil, errors.Wrapf(err, "could not decode mapstructure for '%s'", identifier)
	}

	return val, nil
}

// BuildMetaResource will construct a resource from the given inputs.
func BuildMetaResource(identifier string, data interface{}) (definitions.Resource, error) {
	// Duplicate names are prohibited.
	m := map[string]definitions.Meta{
		noop.Meta{}.ID(): noop.Meta{},
	}

	val, ok := m[identifier]
	if !ok {
		return nil, fmt.Errorf("source '%s' is unknown", identifier)
	}

	if err := mapstructure.Decode(data, &val); err != nil {
		return nil, errors.Wrapf(err, "could not decode mapstructure for '%s'", identifier)
	}

	return val, nil
}

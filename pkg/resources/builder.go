package resources

import (
	"fmt"

	"github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// BuildInputResource will construct a resource from the given inputs.
func BuildInputResource(identifier string, data interface{}, m map[string]definitions.Input) (definitions.Input, error) {
	val, ok := m[identifier]
	if !ok {
		return nil, fmt.Errorf("input resource '%s' is unknown", identifier)
	}

	if err := mapstructure.Decode(data, &val); err != nil {
		return nil, errors.Wrapf(err, "could not decode mapstructure for '%s'", identifier)
	}

	return val, nil
}

// BuildProcessResource will construct a resource from the given process.
func BuildProcessResource(identifier string, data interface{}, m map[string]definitions.Process) (definitions.Process, error) {
	val, ok := m[identifier]
	if !ok {
		return nil, fmt.Errorf("process resource '%s' is unknown", identifier)
	}

	if err := mapstructure.Decode(data, &val); err != nil {
		return nil, errors.Wrapf(err, "could not decode mapstructure for '%s'", identifier)
	}

	return val, nil
}

// BuildOutputResource will construct a resource from the given outputs.
func BuildOutputResource(identifier string, data interface{}, m map[string]definitions.Output) (definitions.Output, error) {
	val, ok := m[identifier]
	if !ok {
		return nil, fmt.Errorf("output resource '%s' is unknown", identifier)
	}

	if err := mapstructure.Decode(data, &val); err != nil {
		return nil, errors.Wrapf(err, "could not decode mapstructure for '%s'", identifier)
	}

	return val, nil
}

// BuildMetaResource will construct a resource from the given inputs.
func BuildMetaResource(identifier string, data interface{}, m map[string]definitions.Meta) (definitions.Meta, error) {
	val, ok := m[identifier]
	if !ok {
		return nil, fmt.Errorf("meta resource '%s' is unknown", identifier)
	}

	if err := mapstructure.Decode(data, &val); err != nil {
		return nil, errors.Wrapf(err, "could not decode mapstructure for '%s'", identifier)
	}

	return val, nil
}

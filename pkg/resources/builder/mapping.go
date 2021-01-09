package builder

import (
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/andrewneudegg/delta/pkg/resources/examples/noop"
	console1 "github.com/andrewneudegg/delta/pkg/resources/utilities/console/v1"
	generators1 "github.com/andrewneudegg/delta/pkg/resources/utilities/generators/v1"
)

// InputMapping defines which resources are what.
func InputMapping() map[string]definitions.Input {
	return map[string]definitions.Input{
		noop.Input{}.ID():        noop.Input{},
		generators1.Input{}.ID(): generators1.Input{},
		console1.Input{}.ID():    console1.Input{},
	}
}

// ProcessMapping defines which resources are what.
func ProcessMapping() map[string]definitions.Process {
	return map[string]definitions.Process{
		noop.Process{}.ID():        noop.Process{},
		generators1.Process{}.ID(): generators1.Process{},
		console1.Process{}.ID():    console1.Process{},
	}
}

// OutputMapping defines which resources are what.
func OutputMapping() map[string]definitions.Output {
	return map[string]definitions.Output{
		noop.Output{}.ID():        noop.Output{},
		generators1.Output{}.ID(): generators1.Output{},
		console1.Output{}.ID():    console1.Output{},
	}
}

// MetaMapping defines which resources are what.
func MetaMapping() map[string]definitions.Meta {
	return map[string]definitions.Meta{
		noop.Meta{}.ID():        noop.Meta{},
		generators1.Meta{}.ID(): generators1.Meta{},
		console1.Meta{}.ID():    console1.Meta{},
	}
}

package builder

import (
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/andrewneudegg/delta/pkg/resources/examples/noop"
)

// InputMapping defines which resources are what.
func InputMapping() map[string]definitions.Input {
	return map[string]definitions.Input{
		noop.Input{}.ID(): noop.Input{},
	}
}

// ProcessMapping defines which resources are what.
func ProcessMapping() map[string]definitions.Process {
	return map[string]definitions.Process{
		noop.Process{}.ID(): noop.Process{},
	}
}

// OutputMapping defines which resources are what.
func OutputMapping() map[string]definitions.Output {
	return map[string]definitions.Output{
		noop.Output{}.ID(): noop.Output{},
	}
}

// MetaMapping defines which resources are what.
func MetaMapping() map[string]definitions.Meta {
	return map[string]definitions.Meta{
		noop.Meta{}.ID(): noop.Meta{},
	}
}

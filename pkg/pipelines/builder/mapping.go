package builder

import (
	"github.com/andrewneudegg/delta/pkg/pipelines/definitions"
	"github.com/andrewneudegg/delta/pkg/pipelines/fipfo"
)

// PipelineMapping returns the mapping between pipeline names and object.
func PipelineMapping() map[string]definitions.Pipeline {
	// Add your Pipeline here.
	return map[string]definitions.Pipeline{
		fipfo.Pipeline{}.ID(): fipfo.Pipeline{},
	}
}

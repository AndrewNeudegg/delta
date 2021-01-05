package fipfo

import (
	"context"
	"fmt"

	"github.com/andrewneudegg/delta/pkg/pipelines/definitions"
)

// Pipeline is a type of pipeline that fans in inputs,
// applies processing, then fans out to the various outputs.
type Pipeline struct {
	Inputs    definitions.PipelineNode `mapstructure:"id"`
	Processes definitions.PipelineNode `mapstructure:"id"`
	Outputs   definitions.PipelineNode `mapstructure:"id"`
}

func (f Pipeline) ID() string {
	return "pipelines/fipfo"
}

func (f Pipeline) Do(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

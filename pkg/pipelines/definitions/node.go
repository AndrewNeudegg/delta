package definitions

// PipelineNode represents a block in a pipeline.
type PipelineNode struct {
	ID     string                 `mapstructure:"id"`     // ID is the name of this thing.
	Config map[string]interface{} `mapstructure:"config"` // Config of this node.
}

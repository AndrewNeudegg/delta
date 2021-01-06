package definitions

// PipelineNode represents a block in a pipeline.
type PipelineNode struct {
	ID     string                 `yaml:"id"`
	Config map[string]interface{} `yaml:"config"`
	Nodes  []PipelineNode         `yaml:"nodes"`
}

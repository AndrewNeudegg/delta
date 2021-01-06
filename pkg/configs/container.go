package configs

import (
	"github.com/andrewneudegg/delta/pkg/pipelines/definitions"
)

// AppSettings is application level configuration.
type AppSettings struct{}

// Container wraps all subconfiguration
type Container struct {
	// ApplicationSettings apply to the whole app.
	ApplicationSettings AppSettings `yaml:"app"`
	// SourceConfigs defines the behaviours of sinks/bridges.
	Pipeline []definitions.PipelineNode `yaml:"pipeline"`
}

// FromBytes will load a container from bytes.
func FromBytes(b []byte) (Container, error) {
	return Container{}, nil
}

// FromFile will load a container from file.
func FromFile(s string) (Container, error) {
	return Container{}, nil
}

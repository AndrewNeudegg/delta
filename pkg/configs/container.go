package configs

import (
	"io/ioutil"

	"github.com/andrewneudegg/delta/pkg/pipelines/definitions"
	"gopkg.in/yaml.v2"
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
	container := Container{}

	err := yaml.Unmarshal(b, &container)
	if err != nil {
		return Container{}, err
	}

	return container, nil
}

// FromFile will load a container from file.
func FromFile(s string) (Container, error) {
	data, err := ioutil.ReadFile(s)
	if err != nil {
		return Container{}, err
	}

	return FromBytes(data)
}

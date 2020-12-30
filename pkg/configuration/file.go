package configuration

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// FileConfig loads configuration from disk.
type FileConfig struct {
	Source string
}

// Load said configuration from disk/
func (f FileConfig) Load() (Container, error) {
	container := Container{}
	data, err := ioutil.ReadFile(f.Source)
	if err != nil {
		return Container{}, err
	}
	err = yaml.Unmarshal(data, &container)
	if err != nil {
		return Container{}, err
	}

	return container, nil
}

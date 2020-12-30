package configuration

import "gopkg.in/yaml.v2"

type RawConfig struct {
	ConfigData []byte
}

func (c RawConfig) Load() (Container, error) {
	container := Container{}

	err := yaml.Unmarshal(c.ConfigData, &container)
	if err != nil {
		return Container{}, err
	}

	return container, nil
}

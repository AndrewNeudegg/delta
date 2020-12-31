package pipeline

import (
	"testing"

	"github.com/andrewneudegg/delta/pkg/configuration"
	"github.com/stretchr/testify/assert"
)

func TestPipelineSmoke(t *testing.T) {

	_, err := BuildPipeline(configuration.Container{
		ApplicationSettings: configuration.AppSettings{},
		SourceConfigs: []configuration.NodeConfig{
			{
				Name: "http/simple",
				Config: map[string]interface{}{
					"ListenAddr":  ":5080",
					"MaxBodySize": 512,
				},
			},
		},
		RelayConfigs: []configuration.NodeConfig{
			{
				Name:   "memory",
				Config: map[string]interface{}{},
			},
		},
		DistributorConfigs: []configuration.NodeConfig{
			{
				Name: "http/direct",
				Config: map[string]interface{}{
					"Addr": "5081",
				},
			},
		},
	})

	assert.Nil(t, err)
}

func TestPipelineSmokeYaml(t *testing.T) {
	config := []byte(`
applicationSettings: {}
sourceConfigurations:
  - name: http/simple
    config:
      ListenAddr: :5080
      MaxBodySize: 512
relayConfigs:
  - name: memory
    config: {}
distributorConfigurations:
  - name: http/direct
    config:
      Addr: http://localhost:5080
`)
	c := configuration.RawConfig{
		ConfigData: config,
	}
	configContainer, err := c.Load()
	assert.Nil(t, err)
	_, err = BuildPipeline(configContainer)
	assert.Nil(t, err)
}

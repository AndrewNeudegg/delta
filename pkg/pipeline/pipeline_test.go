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
				Name: "sink/http",
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
				Name: "naive",
				Config: map[string]interface{}{
					"Addr": "5081",
				},
			},
		},
	})

	assert.Nil(t, err)
}

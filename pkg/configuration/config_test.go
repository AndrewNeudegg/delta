package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationSimple(t *testing.T) {
	config := []byte(`
applicationSettings: {}
sourceConfigurations:
  - name: http/simple
    config:
      ListenAddr: :5050
      MaxBodySize: 512
relayConfigs:
  - name: memory
    config: {}
distributorConfigurations:
  - name: http/direct
    config:
      Addr: http://localhost:5051
`)
	c := RawConfig{
		ConfigData: config,
	}
	lC, err := c.Load()
	assert.Nil(t, err)

	assert.Equal(t, lC.SourceConfigs[0].Name, "http/simple")
	assert.Equal(t, lC.RelayConfigs[0].Name, "memory")
	assert.Equal(t, lC.DistributorConfigs[0].Name, "http/direct")
}

func TestConfigurationSubConfigs(t *testing.T) {
	config := []byte(`
applicationSettings: {}
sourceConfigurations:
  - name: mock/config
    config: {}
    subConfigs: 
    - name: example/example1
      config:
        ListenAddr: :5050
        MaxBodySize: 512
    - name: example/example2
      config:
        ListenAddr: :5050
        MaxBodySize: 512
relayConfigs:
  - name: memory
    config: {}
distributorConfigurations:
  - name: http/direct
    config:
      Addr: http://localhost:5051
`)
	c := RawConfig{
		ConfigData: config,
	}
	lC, err := c.Load()
	assert.Nil(t, err)

	assert.Equal(t, lC.SourceConfigs[0].Name, "mock/config")
	assert.Equal(t, lC.RelayConfigs[0].Name, "memory")
	assert.Equal(t, lC.DistributorConfigs[0].Name, "http/direct")

	assert.Equal(t, lC.SourceConfigs[0].SubConfigs[0].Name, "example/example1")
	assert.Equal(t, lC.SourceConfigs[0].SubConfigs[1].Name, "example/example2")
}

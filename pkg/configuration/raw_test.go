package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawExampleConfig = []byte(`
applicationSettings: {}

# sourceConfigurations
sourceConfigurations:
- name: first
  sourceConfig:
    A: example
    B: example
- name: second
  sourceConfig:
    A: example
    B: example

# relayConfigs
relayConfigs:
- name: thing
  sourceConfig:
    A: example
    B: example

# distributorConfigurations
distributorConfigurations:
- name: thing
  sourceConfig:
    A: example
    B: example
`)

func TestRawSmoke(t *testing.T) {
	f := RawConfig{
		ConfigData: rawExampleConfig,
	}

	_, err := f.Load()
	assert.Nil(t, err)
}

func TestRawBadData(t *testing.T) {
	testBytes := append([]byte("testing"), rawExampleConfig...)
	f := RawConfig{
		ConfigData: testBytes,
	}

	_, err := f.Load()
	assert.Error(t, err)
}

package configuration

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var exampleConfig = []byte(`
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

func TestFileSmoke(t *testing.T) {
	// Setup
	dir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	tmpfn := filepath.Join(dir, "config.yaml")
	if err := ioutil.WriteFile(tmpfn, exampleConfig, 0666); err != nil {
		t.Fatal(err)
	}

	// -- Test --
	fConfig := FileConfig{
		Source: tmpfn,
	}

	config, err := fConfig.Load()
	if err != nil {
		t.Fatal(err)
	}

	sourceConfigs := config.SourceConfigs
	assert.NotEqual(t, 0, len(sourceConfigs))
	assert.Equal(t, "first", sourceConfigs[0].Name)
	assert.Equal(t, "second", sourceConfigs[1].Name)
}

func TestBadFile(t *testing.T) {
	// Setup
	fConfig := FileConfig{
		Source: "doesnt-exist",
	}

	_, err := fConfig.Load()
	assert.Error(t, err)
}

func TestBadFileContent(t *testing.T) {
	testBytes := append([]byte("testing"), exampleConfig...)

	// Setup
	dir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	tmpfn := filepath.Join(dir, "config.yaml")
	if err := ioutil.WriteFile(tmpfn, testBytes, 0666); err != nil {
		t.Fatal(err)
	}

	// -- Test --
	fConfig := FileConfig{
		Source: tmpfn,
	}

	_, err = fConfig.Load()
	assert.Error(t, err)
}

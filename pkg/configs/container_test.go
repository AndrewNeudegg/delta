package configs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	config = []byte(`
applicationSettings: {}
pipeline:
  - id: my/first/pipeline
    config: {}
    nodes:
      - id: nest1/1
        config: {}
        nodes: 
          - id: nest2/1
            config: {}
            nodes: []
          - id: nest2/2
            config: {}
            nodes: []
      - id: nest1/2
        config: {}
        nodes: 
          - id: nest2/1
            config: {}
            nodes: []`)
)

func TestSunshineRaw(t *testing.T) {
	c, err := FromBytes(config)
	assert.Nil(t, err)
	assert.Len(t, c.Pipeline, 1)

	assert.Equal(t, "my/first/pipeline", c.Pipeline[0].ID)
	assert.Len(t, c.Pipeline[0].Nodes, 2)

	assert.Equal(t, "nest1/1", c.Pipeline[0].Nodes[0].ID)
	assert.Equal(t, "nest1/2", c.Pipeline[0].Nodes[1].ID)

	assert.Len(t, c.Pipeline[0].Nodes[0].Nodes, 2)
	assert.Len(t, c.Pipeline[0].Nodes[1].Nodes, 1)

	assert.Equal(t, "nest2/1", c.Pipeline[0].Nodes[0].Nodes[0].ID)
	assert.Equal(t, "nest2/2", c.Pipeline[0].Nodes[0].Nodes[1].ID)
	assert.Equal(t, "nest2/1", c.Pipeline[0].Nodes[1].Nodes[0].ID)
}

func TestSunshineFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	tmpfn := filepath.Join(dir, "config.yaml")
	if err := ioutil.WriteFile(tmpfn, config, 0666); err != nil {
		t.Fatal(err)
	}

	c, err := FromFile(tmpfn)
	assert.Nil(t, err)
	assert.Len(t, c.Pipeline, 1)

	assert.Equal(t, "my/first/pipeline", c.Pipeline[0].ID)
	assert.Len(t, c.Pipeline[0].Nodes, 2)

	assert.Equal(t, "nest1/1", c.Pipeline[0].Nodes[0].ID)
	assert.Equal(t, "nest1/2", c.Pipeline[0].Nodes[1].ID)

	assert.Len(t, c.Pipeline[0].Nodes[0].Nodes, 2)
	assert.Len(t, c.Pipeline[0].Nodes[1].Nodes, 1)

	assert.Equal(t, "nest2/1", c.Pipeline[0].Nodes[0].Nodes[0].ID)
	assert.Equal(t, "nest2/2", c.Pipeline[0].Nodes[0].Nodes[1].ID)
	assert.Equal(t, "nest2/1", c.Pipeline[0].Nodes[1].Nodes[0].ID)
}

func TestRawBadData(t *testing.T) {
	config := append([]byte("testing"), config...)
	_, err := FromBytes(config)
	assert.Error(t, err)
}

func TestBadFilePath(t *testing.T) {
	_, err := FromFile("doesn't exist")
	assert.Error(t, err)
}
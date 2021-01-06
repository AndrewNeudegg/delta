package fipfo

import (
	"context"

	"github.com/pkg/errors"

	"testing"
	"time"

	configuration "github.com/andrewneudegg/delta/pkg/configs"
	"github.com/andrewneudegg/delta/pkg/pipelines"
	"github.com/andrewneudegg/delta/pkg/pipelines/definitions"
	"github.com/stretchr/testify/assert"
)

var (
	config = []byte(`
applicationSettings: {}
pipeline:
  - id: pipelines/fipfo
    config:
      input:
      - id: nest1/1
        config: {}
        nodes: 
          - id: nest2/1
            config: {}
            nodes: []
          - id: nest2/2
            config: {}
            nodes: []
      process:
      - id: nest1/1
        config: {}
        nodes:
          - id: nest2/1
            config: {}
            nodes: []
          - id: nest2/2
            config: {}
            nodes: []
      output:
      - id: nest1/1
        config: {}
        nodes:
          - id: nest2/1
            config: {}
            nodes: []
          - id: nest2/2
            config: {}
            nodes: []`)
)

func TestSunshineBuild(t *testing.T) {
	c, err := configuration.FromBytes(config)
	assert.Nil(t, err)

	_, err = pipelines.BuildPipeline(c.Pipeline[0].ID, c.Pipeline[0].Config, map[string]definitions.Pipeline{
		Pipeline{}.ID(): Pipeline{},
	})
	assert.Nil(t, err)
}

func TestBuildRun(t *testing.T) {
	c, err := configuration.FromBytes(config)
	assert.Nil(t, err)

	p, err := pipelines.BuildPipeline(c.Pipeline[0].ID, c.Pipeline[0].Config, map[string]definitions.Pipeline{
		Pipeline{}.ID(): Pipeline{},
	})
	assert.Nil(t, err)

	go func() {
		err := p.Do(context.TODO())
		if err != nil {
			// panic here because we are in a different goroutine.
			panic(errors.Wrapf(err, "could not Do() for '%s'", p.ID()))
		}
	}()

	time.Sleep(time.Second * 5)
}

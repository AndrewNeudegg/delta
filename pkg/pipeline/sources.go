package pipeline

import (
	"fmt"
	"strings"

	"github.com/andrewneudegg/delta/pkg/configuration"
	metab "github.com/andrewneudegg/delta/pkg/meta/builder"
	"github.com/andrewneudegg/delta/pkg/source"
	sourceb "github.com/andrewneudegg/delta/pkg/source/builder"
	"github.com/pkg/errors"
)

func (p Pipeline) buildS(nc configuration.NodeConfig) (source.S, error) {
	return sourceb.Get(nc.Name, nc.Config)
}

func (p Pipeline) loopS(nc []configuration.NodeConfig) ([]source.S, error) {
	resultantSources := make([]source.S, 0)
	for _, src := range nc {
		if strings.HasPrefix(src.Name, "/meta/") {

			if src.SubConfigs == nil {
				return resultantSources, fmt.Errorf("meta resource '%s' specified but with no sub resources", src.Name)
			}

			m, err := metab.Get(src.Name, src.Config)
			if err != nil {
				return resultantSources, errors.Wrapf(err, "could not build meta for '%s'", src.Name)
			}

			newRelays, err := p.loopS(src.SubConfigs)
			if err != nil {
				return resultantSources, errors.Wrapf(err, "could not build subresources for '%s'", src.Name)
			}

			s, err := m.S(newRelays)
			if err != nil {
				return resultantSources, errors.Wrapf(err, "could not build source for '%s'", src.Name)
			}
			resultantSources = append(resultantSources, s)

		} else {
			s, err := p.buildS(src)
			if err != nil {
				return resultantSources, errors.Wrapf(err, "could not build source for '%s'", src.Name)
			}
			resultantSources = append(resultantSources, s)
		}
	}
	return resultantSources, nil
}

func (p *Pipeline) buildSources(nc []configuration.NodeConfig) error {

	resultantSources, err := p.loopS(nc)
	if err != nil {
		return err
	}

	if len(resultantSources) == 0 {
		return fmt.Errorf("no sources were found")
	}

	p.sources = resultantSources

	return nil
}

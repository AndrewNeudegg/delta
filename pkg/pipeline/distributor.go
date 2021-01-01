package pipeline

import (
	"fmt"
	"strings"

	"github.com/andrewneudegg/delta/pkg/configuration"
	"github.com/andrewneudegg/delta/pkg/distributor"
	distributorb "github.com/andrewneudegg/delta/pkg/distributor/builder"
	metab "github.com/andrewneudegg/delta/pkg/meta/builder"
	"github.com/pkg/errors"
)

func (p Pipeline) buildD(nc configuration.NodeConfig) (distributor.D, error) {
	return distributorb.Get(nc.Name, nc.Config)
}

func (p Pipeline) loopD(nc []configuration.NodeConfig) ([]distributor.D, error) {
	resultantSources := make([]distributor.D, 0)
	for _, src := range nc {
		if strings.HasPrefix(src.Name, "/meta/") {

			if src.SubConfigs == nil {
				return resultantSources, fmt.Errorf("meta resource '%s' specified but with no sub resources", src.Name)
			}

			m, err := metab.Get(src.Name, src.Config)
			if err != nil {
				return resultantSources, errors.Wrapf(err, "could not build meta for '%s'", src.Name)
			}

			newRelays, err := p.loopD(src.SubConfigs)
			if err != nil {
				return resultantSources, errors.Wrapf(err, "could not build subresources for '%s'", src.Name)
			}

			d, err := m.D(newRelays)
			if err != nil {
				return resultantSources, errors.Wrapf(err, "could not build source for '%s'", src.Name)
			}
			resultantSources = append(resultantSources, d)

		} else {
			d, err := p.buildD(src)
			if err != nil {
				return resultantSources, errors.Wrapf(err, "could not build source for '%s'", src.Name)
			}
			resultantSources = append(resultantSources, d)
		}
	}
	return resultantSources, nil
}

func (p *Pipeline) buildDistributors(nc []configuration.NodeConfig) error {

	resultantSources, err := p.loopD(nc)
	if err != nil {
		return err
	}

	if len(resultantSources) == 0 {
		return fmt.Errorf("no sources were found")
	}

	p.distributors = resultantSources

	return nil
}

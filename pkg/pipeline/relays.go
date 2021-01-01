package pipeline

import (
	"fmt"
	"strings"

	"github.com/andrewneudegg/delta/pkg/configuration"
	metab "github.com/andrewneudegg/delta/pkg/meta/builder"
	"github.com/andrewneudegg/delta/pkg/relay"
	relayb "github.com/andrewneudegg/delta/pkg/relay/builder"
	"github.com/pkg/errors"
)

func (p Pipeline) buildR(nc configuration.NodeConfig) (relay.R, error) {
	return relayb.Get(nc.Name, nc.Config)
}

func (p Pipeline) loopR(nc []configuration.NodeConfig) ([]relay.R, error) {
	resultantrelays := make([]relay.R, 0)
	for _, src := range nc {
		if strings.HasPrefix(src.Name, "meta/") {

			if src.SubConfigs == nil {
				return resultantrelays, fmt.Errorf("meta rerelay '%s' specified but with no sub rerelays", src.Name)
			}

			m, err := metab.Get(src.Name, src.Config)
			if err != nil {
				return resultantrelays, errors.Wrapf(err, "could not build meta for '%s'", src.Name)
			}

			newSrcs, err := p.loopR(src.SubConfigs)
			if err != nil {
				return resultantrelays, errors.Wrapf(err, "could not build subrerelays for '%s'", src.Name)
			}

			r, err := m.R(newSrcs)
			if err != nil {
				return resultantrelays, errors.Wrapf(err, "could not build relay for '%s'", src.Name)
			}
			resultantrelays = append(resultantrelays, r)

		} else {
			r, err := p.buildR(src)
			if err != nil {
				return resultantrelays, errors.Wrapf(err, "could not build relay for '%s'", src.Name)
			}
			resultantrelays = append(resultantrelays, r)
		}
	}
	return resultantrelays, nil
}

func (p *Pipeline) buildRelays(nc []configuration.NodeConfig) error {

	resultantrelays, err := p.loopR(nc)
	if err != nil {
		return err
	}

	if len(resultantrelays) == 0 {
		return fmt.Errorf("no relays were found")
	}

	p.relays = resultantrelays

	return nil
}

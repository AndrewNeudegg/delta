package pipeline

import (
	"github.com/andrewneudegg/delta/pkg/configuration"
	"github.com/andrewneudegg/delta/pkg/distributor"
	distributorb "github.com/andrewneudegg/delta/pkg/distributor/builder"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
	relayb "github.com/andrewneudegg/delta/pkg/relay/builder"
	"github.com/andrewneudegg/delta/pkg/source"
	sourceb "github.com/andrewneudegg/delta/pkg/source/builder"
)

type Pipeline struct {
	inCh  chan events.Event // inCh merges all input channels into a singular channel.
	outCh chan events.Event // outCh replicates events to all output channels.

	sources      []source.S
	relays       []relay.R
	distributors []distributor.D
}

func BuildPipeline(c configuration.Container) (Pipeline, error) {
	p := Pipeline{
		inCh:         make(chan events.Event),
		outCh:        make(chan events.Event),
		sources:      make([]source.S, 0),
		distributors: make([]distributor.D, 0),
	}

	for _, sConfig := range c.SourceConfigs {
		s, err := sourceb.Get(sConfig.Name, sConfig.Config)
		if err != nil {
			return Pipeline{}, err
		}
		p.sources = append(p.sources, s)
	}

	for _, rConfig := range c.RelayConfigs {
		r, err := relayb.Get(rConfig.Name, rConfig.Config)
		if err != nil {
			return Pipeline{}, err
		}
		p.relays = append(p.relays, r)
	}

	for _, dConfig := range c.DistributorConfigs {
		d, err := distributorb.Get(dConfig.Name, dConfig.Config)
		if err != nil {
			return Pipeline{}, err
		}
		p.distributors = append(p.distributors, d)
	}

	// Now we have constructed each of the nodes we must connect them.

	return p, nil
}

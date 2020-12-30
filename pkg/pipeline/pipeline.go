package pipeline

import (
	"context"

	"github.com/andrewneudegg/delta/pkg/configuration"
	"github.com/andrewneudegg/delta/pkg/distributor"
	distributorb "github.com/andrewneudegg/delta/pkg/distributor/builder"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
	relayb "github.com/andrewneudegg/delta/pkg/relay/builder"
	"github.com/andrewneudegg/delta/pkg/source"
	sourceb "github.com/andrewneudegg/delta/pkg/source/builder"
	log "github.com/sirupsen/logrus"
)

// Pipeline is the representation of data flow through this application.
type Pipeline struct {
	inCh  chan events.Event // inCh merges all input channels into a singular channel.
	outCh chan events.Event // outCh replicates events to all output channels.

	sources      []source.S
	relays       []relay.R
	distributors []distributor.D
}

// BuildPipeline will construct the pipeline at the core of delta.
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
	// If we want to insert middleware, i.e. telemetry, this is
	// probably the place to do it.
	sourceChannels := []chan events.Event{}
	for _, s := range p.sources {
		thisSourceChan := make(chan events.Event)
		go func(s source.S, ch chan events.Event) {
			err := s.Do(context.TODO(), ch)
			log.Error(err)
		}(s, thisSourceChan)
		sourceChannels = append(sourceChannels, thisSourceChan)
	}

	go func() {
		fanIn(context.TODO(), sourceChannels, p.inCh)
	}()

	// --  --

	var previousSourceOutput *chan events.Event
	previousSourceOutput = &p.inCh
	for _, r := range p.relays {
		thisRelayOutputChan := make(chan events.Event)
		go r.Do(context.TODO(), *previousSourceOutput, thisRelayOutputChan)
		previousSourceOutput = &thisRelayOutputChan
	}

	distributorChannels := []chan events.Event{}
	for _, d := range p.distributors {
		distributorInputChannel := make(chan events.Event)
		go func(d distributor.D, ch chan events.Event) {
			err := d.Do(context.TODO(), ch)
			log.Error(err)
		}(d, distributorInputChannel)
		distributorChannels = append(distributorChannels, distributorInputChannel)
	}

	go func() {
		fanOut(context.TODO(), *previousSourceOutput, distributorChannels)
	}()

	return p, nil
}

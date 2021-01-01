package pipeline

import (
	"context"
	"fmt"
	"sync"

	"github.com/andrewneudegg/delta/pkg/configuration"
	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
	"github.com/andrewneudegg/delta/pkg/source"
	log "github.com/sirupsen/logrus"
)

// Pipeline is the representation of data flow through this application.
type Pipeline struct {
	sources      []source.S
	relays       []relay.R
	distributors []distributor.D
}

// Await blocks indefinitely.
func (p Pipeline) Await() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

// BuildPipeline will construct the pipeline at the core of delta.
func BuildPipeline(c configuration.Container) (Pipeline, error) {
	p := Pipeline{
		sources:      make([]source.S, 0),
		distributors: make([]distributor.D, 0),
		relays:       make([]relay.R, 0),
	}

	if err := p.buildSources(c.SourceConfigs); err != nil {
		return Pipeline{}, err
	}

	if err := p.buildRelays(c.RelayConfigs); err != nil {
		return Pipeline{}, err
	}

	if err := p.buildDistributors(c.DistributorConfigs); err != nil {
		return Pipeline{}, err
	}

	// Now we have constructed each of the nodes we must connect them.
	// If we want to insert middleware, i.e. telemetry, this is
	// probably the place to do it.
	sourceChannels := []chan events.Event{}
	for _, s := range p.sources {
		if s == nil {
			return Pipeline{}, fmt.Errorf("source was unexpectedly nil")
		}

		thisSourceChan := make(chan events.Event)
		go func(s source.S, ch chan events.Event) {
			err := s.Do(context.TODO(), ch)
			log.Error(err)
		}(s, thisSourceChan)
		sourceChannels = append(sourceChannels, thisSourceChan)
	}

	// TODO: Use prometheus middleware here instead...
	inCh := Inject(make(chan events.Event), NoopEventMiddleware("inbound"))
	go func() {
		fanIn(context.TODO(), sourceChannels, inCh)
	}()

	// --  --

	var previousSourceOutput *chan events.Event
	previousSourceOutput = &inCh
	for _, r := range p.relays {

		if r == nil {
			return Pipeline{}, fmt.Errorf("relay was unexpectedly nil")
		}

		thisRelayOutputChan := make(chan events.Event)
		go r.Do(context.TODO(), *previousSourceOutput, thisRelayOutputChan)
		previousSourceOutput = &thisRelayOutputChan
	}

	distributorChannels := []chan events.Event{}
	for _, d := range p.distributors {

		if d == nil {
			return Pipeline{}, fmt.Errorf("distributor was unexpectedly nil")
		}

		distributorInputChannel := make(chan events.Event)
		go func(d distributor.D, ch chan events.Event) {
			err := d.Do(context.TODO(), ch)
			log.Error(err)
		}(d, distributorInputChannel)
		distributorChannels = append(distributorChannels, distributorInputChannel)
	}

	go func() {
		// TODO: Use prometheus middleware here instead...
		fanOut(context.TODO(), Inject(*previousSourceOutput, NoopEventMiddleware("outbound")), distributorChannels)
	}()

	return p, nil
}

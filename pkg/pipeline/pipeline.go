package pipeline

import (
	"context"
	"sync"

	"github.com/andrewneudegg/delta/pkg/configuration"
	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
	"github.com/andrewneudegg/delta/pkg/source"
	"github.com/andrewneudegg/delta/pkg/utils"
	"github.com/pkg/errors"
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

// hookupSources will build a channel for each source so that
// they may be audited at another point in time.
func (p *Pipeline) hookupSources(ctx context.Context) (chan []events.Event, error) {
	sourceChannels := []chan []events.Event{}
	outputCh := make(chan []events.Event)

	for _, s := range p.sources {
		thisSourceChan := make(chan []events.Event)

		go func(s source.S, ch chan []events.Event) {
			log.Infof("launching source '%s'", s.ID())
			err := s.SDo(context.TODO(), ch)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to do source"))
			}
			log.Warnf("source '%s' has exited", s.ID())
		}(s, thisSourceChan)
		sourceChannels = append(sourceChannels, thisSourceChan)
	}
	log.Infof("found '%d' sources", len(p.sources))

	c := utils.Channels{}
	go c.FanIn(ctx, sourceChannels, outputCh)

	return outputCh, nil
}

// hookupRelays will daisy chain the relays together
// from first to last, returning the final output chan.
func (p *Pipeline) hookupRelays(ctx context.Context, input chan []events.Event) (chan []events.Event, error) {
	var previousSourceOutput *chan []events.Event
	previousSourceOutput = &input

	for _, r := range p.relays {
		thisRelayOutputChan := make(chan []events.Event)
		go func(inCh <-chan []events.Event, outCh chan<- []events.Event) {
			log.Infof("launching relay '%s'", r.ID())
			err := r.RDo(context.TODO(), inCh, outCh)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to do relay"))
			}
			log.Warnf("relay '%s' has exited", r.ID())
		}(*previousSourceOutput, thisRelayOutputChan)
		previousSourceOutput = &thisRelayOutputChan
	}
	log.Infof("found '%d' relays", len(p.relays))

	return *previousSourceOutput, nil
}

func (p *Pipeline) hookupDistributors(ctx context.Context, input chan []events.Event) error {
	distributorChannels := make([]chan []events.Event, len(p.distributors))
	for _, d := range p.distributors {
		distributorInputChannel := make(chan []events.Event)
		go func(d distributor.D, ch chan []events.Event) {
			log.Infof("launching distributor '%s'", d.ID())
			err := d.DDo(context.TODO(), ch)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to do distributor"))
			}
			log.Warnf("distributor '%s' has exited", d.ID())
		}(d, distributorInputChannel)
		distributorChannels = append(distributorChannels, distributorInputChannel)
	}
	log.Infof("found '%d' distributors", len(p.distributors))

	c := utils.Channels{}
	go c.FanOut(ctx, input, distributorChannels)

	return nil
}

// Hookup will build the channels required for data flow and press go.
func (p *Pipeline) Hookup(ctx context.Context) error {
	sourceCh, err := p.hookupSources(ctx)
	if err != nil {
		return err
	}

	finalRelayCh, err := p.hookupRelays(ctx, sourceCh)
	if err != nil {
		return err
	}

	return p.hookupDistributors(ctx, finalRelayCh)
}

// BuildPipeline will construct the pipeline at the core of delta.
func BuildPipeline(c configuration.Container) (Pipeline, error) {
	p := Pipeline{
		sources:      make([]source.S, 0),
		distributors: make([]distributor.D, 0),
		relays:       make([]relay.R, 0),
	}

	if err := p.buildSources(c.SourceConfigs); err != nil {
		return Pipeline{}, errors.Wrap(err, "could not build sources")
	}

	if err := p.buildRelays(c.RelayConfigs); err != nil {
		return Pipeline{}, errors.Wrap(err, "could not build relays")
	}

	if err := p.buildDistributors(c.DistributorConfigs); err != nil {
		return Pipeline{}, errors.Wrap(err, "could not build distributors")
	}

	if err := p.Hookup(context.TODO()); err != nil {
		return Pipeline{}, errors.Wrap(err, "could not hookup sources, relays & distributors")
	}

	// Now we have constructed each of the nodes we must connect them.
	// If we want to insert middleware, i.e. telemetry, this is
	// probably the place to do it.

	// TODO: Use prometheus middleware here instead...
	// inCh := Inject(make(chan []events.Event), NoopEventMiddleware("inbound"))
	// go func() {
	// 	chUtils.FanIn(context.TODO(), sourceChannels, inCh)
	// }()

	return p, nil
}

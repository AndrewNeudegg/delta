package chaos

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
	"github.com/andrewneudegg/delta/pkg/source"
	"github.com/andrewneudegg/delta/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ID = "meta/chaos/simple"
)

// Simple intercepts and randomly fails events.
type Simple struct {
	FailChance float32 `mapstructure:"failChance"` // 0.5

	sources      []source.S
	relays       []relay.R
	distributors []distributor.D
}

// ID is a human readable string identifying this thing.
func (m Simple) ID() string {
	return ID
}

// ----------- Substructs -----------
type s struct {
	source.S
	m *Simple
}

// ID is a human readable string identifying this thing.
func (s s) ID() string {
	return ID
}

func (s s) Do(ctx context.Context, ch chan<- events.Event) error {
	return s.m.doS(ctx, ch)
}

type r struct {
	relay.R
	m *Simple
}

// ID is a human readable string identifying this thing.
func (r r) ID() string {
	return ID
}

func (r r) Do(ctx context.Context, outbound <-chan events.Event, inbound chan<- events.Event) error {
	return r.m.doR(ctx, outbound, inbound)
}

type d struct {
	distributor.D
	m *Simple
}

// ID is a human readable string identifying this thing.
func (d d) ID() string {
	return ID
}

func (d d) Do(ctx context.Context, ch <-chan events.Event) error {
	return d.m.doD(ctx, ch)
}

// ----------- Substructs -----------

// S source intermediary.
func (m *Simple) S(sources []source.S) (source.S, error) {
	m.sources = sources
	return s{m: m}, nil
}

// R relay intermediary.
func (m *Simple) R(relays []relay.R) (relay.R, error) {
	m.relays = relays
	return r{m: m}, nil
}

// D distributor intermediary.
func (m *Simple) D(distributors []distributor.D) (distributor.D, error) {
	m.distributors = distributors
	return d{m: m}, nil
}

// ----------- /Substructs -----------

func (m Simple) isChance(f float32) bool {
	// rand.Float64() == 0.1, 0.5, 0.8
	// if f == 0.1 (10% chance) then rand has to be above 0.9.
	// if f == 0.90 (90% chance) then rand has to be above 0.1.
	return rand.Float32() > f
}

// DoS will do S with some modification.
func (m *Simple) doS(ctx context.Context, ch chan<- events.Event) error {
	wg := sync.WaitGroup{}
	wg.Add(1)

	// We don't do much here, we just pass on the do.
	for _, v := range m.sources {
		go v.Do(ctx, ch)
	}
	wg.Wait()

	return ctx.Err()
}

// DoR will do R with some modification.
func (m *Simple) doR(ctx context.Context, outbound <-chan events.Event, inbound chan<- events.Event) error {
	wg := sync.WaitGroup{}
	wg.Add(1)

	// We don't do much here, we just pass on the do.
	for _, v := range m.relays {
		go v.Do(ctx, outbound, inbound)
	}

	wg.Wait()

	return ctx.Err()
}

// DoD will do D with some modification.
func (m *Simple) doD(ctx context.Context, ch <-chan events.Event) error {
	c := utils.Channels{}
	wg := sync.WaitGroup{}
	wg.Add(1)

	chs := make([]chan events.Event, len(m.distributors))
	intermediateCh := make(chan events.Event)
	for i := 0; i < len(m.distributors); i++ {
		chs[i] = make(chan events.Event)
		go func(i int) {
			log.Infof("starting proxy distributor '%s'", m.distributors[i].ID())
			err := m.distributors[i].Do(ctx, chs[i])
			if err != nil {
				log.Error(errors.Wrapf(err, "an error occurred starting, '%s' chaos proxy distributor", m.distributors[i].ID()))
			}
			log.Warnf("proxy distributor '%s' has exited", m.distributors[i].ID())
		}(i)

	}

	go c.FanOut(ctx, intermediateCh, chs)

	go func() {
		for ctx.Err() == nil {
			select {
			case e := <-ch:
				if m.isChance(m.FailChance) {
					// log.Debugf("event '%s' was lucky!", e.GetMessageID())
					intermediateCh <- e
					continue
				} else {
					// log.Debugf("event '%s' was unlucky and will be dropped, event drop probability is '%f'", e.GetMessageID(), m.FailChance)
					e.Fail(fmt.Errorf("chaos happened to this event"))
				}
			case _ = <-ctx.Done():
				break
			}
		}
	}()

	wg.Wait()

	return nil
}

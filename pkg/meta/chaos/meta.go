package chaos

import (
	"context"
	"math/rand"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
	"github.com/andrewneudegg/delta/pkg/source"
)

// Simple intercepts and randomly fails events.
type Simple struct {
	FailChance float32 // 0.5
}

// ----------- Substructs -----------
type s struct {
	source.S
	m *Simple
}

func (s s) Do(ctx context.Context, ch chan<- events.Event) error {
	return s.m.DoS(ctx, ch)
}

type r struct {
	relay.R
	m *Simple
}

func (r r) Do(ctx context.Context, outbound <-chan events.Event, inbound chan<- events.Event) error {
	return r.m.DoR(ctx, outbound, inbound)
}

type d struct {
	distributor.D
	m *Simple
}

func (d d) Do(ctx context.Context, ch <-chan events.Event) error {
	return d.m.DoD(ctx, ch)
}

// ----------- Substructs -----------

// S source intermediary.
func (m Simple) S([]source.S) (source.S, error) { return s{m: &m}, nil }

// R relay intermediary.
func (m Simple) R([]relay.R) (relay.R, error) { return r{m: &m}, nil }

// D distributor intermediary.
func (m Simple) D([]distributor.D) (distributor.D, error) { return d{m: &m}, nil }

func (m Simple) isChance(f float32) bool {
	// rand.Float64() == 0.1, 0.5, 0.8
	// if f == 0.1 (10% chance) then rand has to be above 0.9.
	// if f == 0.90 (90% chance) then rand has to be above 0.1.
	return rand.Float32() > (1 - f)
}

// DoS will do S with some modification.
func (m Simple) DoS(context.Context, chan<- events.Event) error {
	return nil

	// nCh := make(chan events.Event)
	// go func() {
	// 	for {
	// 		select {
	// 		case e := <-ch:
	// 			if m.isChance(m.FailChance) {
	// 				e.Fail(fmt.Errorf("event was unlucky"))
	// 				continue
	// 			}

	// 			// if its lucky then continue...
	// 			nCh <- e
	// 		case _ = <-ctx.Done():
	// 			return
	// 		}
	// 	}
	// }()

	// return s.Do(ctx, nCh)
}

// DoR will do R with some modification.
func (m Simple) DoR(ctx context.Context, outbound <-chan events.Event, inbound chan<- events.Event) error {
	return nil

	// nCh := make(chan events.Event)

	// go func() {
	// 	for {
	// 		select {
	// 		case e := <-chOut:
	// 			if m.isChance(m.FailChance) {
	// 				e.Fail(fmt.Errorf("event was unlucky"))
	// 				continue
	// 			}

	// 			// if its lucky then continue...
	// 			nCh <- e
	// 		case _ = <-ctx.Done():
	// 			return
	// 		}
	// 	}
	// }()

	// return r.Do(ctx, nCh, chIn)
}

// DoD will do D with some modification.
func (m Simple) DoD(ctx context.Context, ch <-chan events.Event) error {
	return nil

	// nCh := make(chan events.Event)

	// go func() {
	// 	for {
	// 		select {
	// 		case e := <-ch:
	// 			if m.isChance(m.FailChance) {
	// 				e.Fail(fmt.Errorf("event was unlucky"))
	// 				continue
	// 			}

	// 			// if its lucky then continue...
	// 			nCh <- e
	// 		case _ = <-ctx.Done():
	// 			return
	// 		}
	// 	}
	// }()

	// return d.Do(ctx, nCh)
}

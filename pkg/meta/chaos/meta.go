package chaos

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/andrewneudegg/delta/pkg/distributor"
	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
	"github.com/andrewneudegg/delta/pkg/source"
)

// ChaosSimple intercepts and randomly fails events.
type ChaosSimple struct {
	FailChance float32 // 0.5
}

func (m ChaosSimple) isChance(f float32) bool {
	// rand.Float64() == 0.1, 0.5, 0.8
	// if f == 0.1 (10% chance) then rand has to be above 0.9.
	// if f == 0.90 (90% chance) then rand has to be above 0.1.
	return rand.Float32() > (1 - f)
}

// DoS will do S with some modification.
func (m ChaosSimple) DoS(ctx context.Context, ch chan events.Event, s source.S) error {
	nCh := make(chan events.Event)
	go func() {
		for {
			select {
			case e := <-ch:
				if m.isChance(m.FailChance) {
					e.Fail(fmt.Errorf("event was unlucky"))
					continue
				}

				// if its lucky then continue...
				nCh <- e
			case _ = <-ctx.Done():
				return
			}
		}
	}()

	return s.Do(ctx, nCh)
}

// DoR will do R with some modification.
func (m ChaosSimple) DoR(ctx context.Context, chOut chan events.Event, chIn chan events.Event, r relay.R) error {
	nCh := make(chan events.Event)

	go func() {
		for {
			select {
			case e := <-chOut:
				if m.isChance(m.FailChance) {
					e.Fail(fmt.Errorf("event was unlucky"))
					continue
				}

				// if its lucky then continue...
				nCh <- e
			case _ = <-ctx.Done():
				return
			}
		}
	}()

	return r.Do(ctx, nCh, chIn)
}

// DoD will do D with some modification.
func (m ChaosSimple) DoD(ctx context.Context, ch chan events.Event, d distributor.D) error {
	nCh := make(chan events.Event)

	go func() {
		for {
			select {
			case e := <-ch:
				if m.isChance(m.FailChance) {
					e.Fail(fmt.Errorf("event was unlucky"))
					continue
				}

				// if its lucky then continue...
				nCh <- e
			case _ = <-ctx.Done():
				return
			}
		}
	}()

	return d.Do(ctx, nCh)
}

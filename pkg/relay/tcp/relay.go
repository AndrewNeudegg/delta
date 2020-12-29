package tcp

import (
	"context"
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
)

// Opts specifies where messages should be routed.
type Opts interface {
	Protocol() string
	SendAddress() string
	ListenAddress() string
}

// Relay passes events from a source to a destination via any intermediary.
type Relay struct {
	relay.R
	Opts Opts
}

// Listen for incoming events.
func (r Relay) Listen(ctx context.Context, ch chan<- events.Event) error {
	handleConnection := func(conn net.Conn) {
		defer conn.Close()
		// This is a little bit wasteful.
		dec := gob.NewDecoder(conn)
		e := &events.EventMsg{}
		dec.Decode(e)
		ch <- e
	}

	ln, err := net.Listen(r.Opts.Protocol(), r.Opts.ListenAddress())
	if err != nil {
		return err
	}

	// func will break if cancellation
	// this is because ln.Accept() will block
	// until a connection or error.
	go func() {
		select {
		case _ = <-ctx.Done():
			ln.Close()
		}
	}()

	for ctx.Err() == nil {
		conn, err := ln.Accept()
		if err != nil {
			conn.Close()
			continue
		}
		go handleConnection(conn)
	}

	return nil
}

// Emit will emit all messages passed in to the channel.
func (r Relay) Emit(ctx context.Context, ch <-chan events.Event) error {

	emit := func(e events.Event) error {
		conn, err := net.Dial(r.Opts.Protocol(), r.Opts.SendAddress())
		defer conn.Close()
		if err != nil {
			return err
		}
		encoder := gob.NewEncoder(conn)
		return encoder.Encode(e)
	}

	backoffRetry := func(e events.Event) error {
		backoffSecs := time.Duration(1) * time.Second
		for i := 0; i < 5; i++ {
			time.Sleep(backoffSecs)
			if err := emit(e); err != nil {
				backoffSecs = backoffSecs * 2
				continue
			} else {
				return nil
			}
		}
		return fmt.Errorf("retry backoff reached")
	}

	// loop, breaking on either ctx done or event msg.
	for ctx.Err() == nil {
		select {
		case e := <-ch:
			go backoffRetry(e)
		case _ = <-ctx.Done():
			break
		}
	}

	return ctx.Err()
}

// Do will begin processing on the I/O channels.
func (r Relay) Do(ctx context.Context, outbound <-chan events.Event, inbound chan<- events.Event) error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	// start the tcp server and begin listening.
	go func() {
		r.Listen(ctx, inbound)
		wg.Done()
	}()

	// start the tcp client and begin emitting.
	go func() {
		r.Emit(ctx, outbound)
		wg.Done()
	}()

	// wait until we exit.
	wg.Done()
	return ctx.Err()
}

package tcp

import (
	"context"
	"encoding/gob"
	"net"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/messaging"
)

// Relay passes event messages downstream.
type Relay struct {
	messaging.Relay

	protocol   string // tcp
	listenAddr string // :8080

	inboundChan chan<- events.Event
}

// listen will start the tcp server and fork a goroutine for each connection.
func (r Relay) listen(ctx context.Context, ch chan<- events.Event) error {

	handleConnection := func(conn net.Conn) {
		defer conn.Close()
		// This is a little bit wasteful.
		dec := gob.NewDecoder(conn)
		e := &events.EventMsg{}
		dec.Decode(e)
		ch <- e
	}

	ln, err := net.Listen(r.protocol, r.listenAddr)
	if err != nil {
		return err
	}

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

// Do will pass events from a channel to a target, by looking them up.
func (r Relay) Do(ctx context.Context, ch <-chan events.Event, p Phonebook) error {
	go r.listen(ctx, r.inboundChan)

	for ctx.Err() == nil {
		e := <-ch
		t, err := p.Lookup(e)
		if err != nil {
			return err
		}
		err = t.Send(ctx, e)
		if err != nil {
			return err
		}
	}

	return ctx.Err()
}

package tcp

import (
	"fmt"
	"math/rand"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/messaging"
)

// Phonebook is a simple impl that will pass messages to the defined channel.
type Phonebook struct {
	messaging.Phonebook

	KnownAddresses []string
}

func (p Phonebook) randAddr() (string, error) {
	if len(p.KnownAddresses) == 0 {
		return "", fmt.Errorf("no addresses were supplied")
	}

	return p.KnownAddresses[rand.Intn(len(p.KnownAddresses))], nil
}

// Lookup will find a target suitable for the given event.
func (p Phonebook) Lookup(e events.Event) (Target, error) {
	addr, err := p.randAddr()
	if err != nil {
		return Target{}, err
	}

	return Target{
		protocol: "tcp",
		address:  addr,
	}, nil
}

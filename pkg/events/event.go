package events

import (
	"encoding/json"
)

// Distributor will queue an event for distribution, taking care of the messy details.
type Distributor interface {
	Enqueue(Event) error // Enqueue a new event.
}

// Event is the most abstract description of the message passing.
type Event interface {
	GetMessageID() *string          // GetMessageID returns the message identifier.
	GetHeaders() *map[string]string // GetHeaders returns HTTP compatible headers.
	GetURI() *string                // GetURI returns a HTTP like route.
	GetContent() *[]byte            // GetContent returns the []byte content of the message.
}

// EventMsg is the underlying Event.
type EventMsg struct {
	ID      string            `json:"id"`
	Headers map[string]string `json:"headers"`
	URI     string            `json:"uri"`
	Content []byte            `json:"content"`
}

// ToJSON will convert an EventMsg to JSON.
func (e *EventMsg) ToJSON() (string, error) {
	jsonB, err := e.ToJSONb()
	return string(jsonB), err
}

// ToJSONb will convert an EventMsg to JSON []Byte.
func (e *EventMsg) ToJSONb() ([]byte, error) {
	return json.Marshal(*e)
}

// GetMessageID returns the message identifier.
func (e *EventMsg) GetMessageID() *string {
	return &e.ID
}

// GetHeaders returns HTTP compatible headers.
func (e *EventMsg) GetHeaders() *map[string]string {
	return &e.Headers
}

// GetURI returns a HTTP like route.
func (e *EventMsg) GetURI() *string {
	return &e.URI
}

// GetContent returns the []byte content of the message.
func (e *EventMsg) GetContent() *[]byte {
	return &e.Content
}

// FromJSON will load an EventMsg from JSON string.
func FromJSON(payload string) (*EventMsg, error) {
	return FromJSONb([]byte(payload))
}

// FromJSONb will load an EventMsg from JSON []Byte.
func FromJSONb(payload []byte) (*EventMsg, error) {
	var e EventMsg
	err := json.Unmarshal(payload, &e)
	return &e, err
}
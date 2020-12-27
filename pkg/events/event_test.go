package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E(t *testing.T) {
	event := EventMsg{
		ID: "test",
		Headers: map[string]string{
			"hello": "world",
		},
		URI:     "/example",
		Content: []byte("hello world"),
	}

	assert.Equal(t, *event.GetMessageID(), event.ID)
	assert.Equal(t, *event.GetHeaders(), event.Headers)
	assert.Equal(t, *event.GetURI(), event.URI)
	assert.Equal(t, *event.GetContent(), event.Content)
}

func TestE2EJson(t *testing.T) {
	event := EventMsg{
		ID: "test",
		Headers: map[string]string{
			"hello": "world",
		},
		URI:     "/example",
		Content: []byte("hello world"),
	}

	jsonEvent, err := event.ToJSON()
	assert.Nil(t, err, "could not get JSON")

	reloadedEvent, err := FromJSON(jsonEvent)
	assert.Nil(t, err, "could not convert back from JSON")

	assert.Equal(t, *reloadedEvent.GetMessageID(), *event.GetMessageID())
	assert.Equal(t, *reloadedEvent.GetHeaders(), *event.GetHeaders())
	assert.Equal(t, *reloadedEvent.GetURI(), *event.GetURI())
	assert.Equal(t, *reloadedEvent.GetContent(), *event.GetContent())
}

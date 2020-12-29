package _queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// These tests should pass regardless of the specific queue implementation.
// Consider it the reference.

type Element struct {
	id string
}

func (e Element) ID() string {
	return e.id
}

func TestSmoke(t *testing.T) {
	q := NewMemoryQ()

	elems := []Element{
		{
			id: "hello",
		},
		{
			id: "world",
		},
		{
			id: "!",
		},
	}

	for _, v := range elems {
		q.Push(v)
	}

	for _, v := range elems {
		result, err := q.Pop()
		assert.Nil(t, err)
		assert.Equal(t, result.(Element).ID(), v.ID())
	}
}

func TestFlow(t *testing.T) {
	q := NewMemoryQ()

	elems := []Element{
		{
			id: "hello",
		},
		{
			id: "world",
		},
	}

	for _, v := range elems {
		q.Push(v)
	}

	result, err := q.Pop()
	assert.Nil(t, err)
	assert.Equal(t, result.(Element).ID(), "hello")
	err = q.Retry(result)
	assert.Nil(t, err)

	result, err = q.Pop()
	assert.Nil(t, err)
	assert.Equal(t, result.(Element).ID(), "world")

	result, err = q.Pop()
	assert.Nil(t, err)
	assert.Equal(t, result.(Element).ID(), "hello")
}

func TestBulkRightN(t *testing.T) {
	q := NewMemoryQ()

	elems := []Element{
		{
			id: "hello",
		},
		{
			id: "world",
		},
		{
			id: "!",
		},
	}

	for _, v := range elems {
		q.Push(v)
	}

	popedElems, err := q.PopN(len(elems))
	assert.Nil(t, err)

	for i := range popedElems {
		assert.Nil(t, err)
		assert.Equal(t, popedElems[i].(Element).ID(), elems[i].ID())
	}
}

func TestBulkBigN(t *testing.T) {
	q := NewMemoryQ()

	elems := []Element{
		{
			id: "hello",
		},
		{
			id: "world",
		},
		{
			id: "!",
		},
	}

	for _, v := range elems {
		q.Push(v)
	}

	popedElems, err := q.PopN(len(elems) + 2)
	assert.Nil(t, err)

	for i := range popedElems {
		assert.Nil(t, err)
		assert.Equal(t, popedElems[i].(Element).ID(), elems[i].ID())
	}
}

func TestLen(t *testing.T) {
	q := NewMemoryQ()

	elems := []Element{
		{
			id: "hello",
		},
		{
			id: "world",
		},
		{
			id: "!",
		},
	}

	for _, v := range elems {
		q.Push(v)
	}

	assert.Equal(t, q.Len(), 3)

	q.Pop()

	assert.Equal(t, q.Len(), 2)

}

func TestFailMsg(t *testing.T) {
	q := NewMemoryQ()

	elems := []Element{
		{
			id: "hello",
		},
		{
			id: "world",
		},
		{
			id: "!",
		},
	}

	for _, v := range elems {
		q.Push(v)
	}

	assert.Equal(t, q.Len(), 3)

	thing, _ := q.Pop()
	q.Failed(thing)

	assert.Equal(t, q.Len(), 2)
}

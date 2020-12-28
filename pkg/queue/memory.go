package queue

import (
	"container/list"
)

// memoryQ is a memory only queue.
type memoryQ struct {
	Q

	list *list.List
}

// Len will return the current number of elements in the queue.
func (q *memoryQ) Len() int {
	return q.list.Len()
}

// Pop will return an element from the front of the queue.
func (q *memoryQ) Pop() (interface{}, error) {
	elem := q.list.Front()
	q.list.Remove(elem)
	return elem.Value, nil
}

// PopN will return up to N elements from the front of the queue.
func (q *memoryQ) PopN(n int) ([]interface{}, error) {
	if q.list.Len() < n {
		n = q.list.Len()
	}
	results := make([]interface{}, n)

	for i := 0; i < n; i++ {
		elem := q.list.Front()
		q.list.Remove(elem)
		results[i] = elem.Value
	}

	return results, nil
}

// Push will add an element to the back of the queue.
func (q *memoryQ) Push(e interface{}) error {
	q.list.PushBack(e)
	return nil
}

// Retry will add an element to the back of the queue.
func (q *memoryQ) Retry(e interface{}) error {
	return q.Push(e)
}

// Failed will do nothing, but other queues may wish to dead letter these messages.
func (q *memoryQ) Failed(e interface{}) error {
	// Do nothing, we may wish to make additional calls here.
	return nil
}

// NewMemoryQ will return a new instance of the memory queue.
func NewMemoryQ() Q {
	return &memoryQ{
		list: list.New(),
	}
}

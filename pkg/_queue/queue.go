package _queue

// Q is a any type of queue, memory, RPC, HTTP, Kafka, etc.
type Q interface {
	Len() int

	Pop() (interface{}, error)
	PopN(int) ([]interface{}, error)
	Push(interface{}) error

	Retry(interface{}) error
	Failed(interface{}) error
}

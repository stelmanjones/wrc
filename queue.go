package wrc

import "errors"

var (
	// ErrQueueEmpty is thrown if trying to access an empty queue.
	ErrQueueEmpty = errors.New("queue empty")

	// ErrQueueFull is thrown if trying to append to a queue that is full.
	ErrQueueFull = errors.New("queue full")
)

// Queue defines a generic FIFO queue interface.
type Queue[T any] interface {
	Enqueue(T) error
	Dequeue() (T, error)
	Resize(s int)
	Len() int
}

// A generic FIFO queue.
type FIFOQueue[T any] struct {
	data []T
	maxSize int
	}

// Creates a new queue of type T.
func NewFIFOQueue[T any]() *FIFOQueue[T] {
	return &FIFOQueue[T]{
		data: make([]T, 1000),
		maxSize: 1000,
	}
}

// Returns the queue length.
func (q *FIFOQueue[T]) Len() int {
	return len(*&q.data)
}

// Enqueue's a new value to the queue.
func (q *FIFOQueue[T]) Enqueue(value T) error {
	q.data = append(q.data, value)
	return nil
}

func (q *FIFOQueue[T]) Data() ([]T, error) {
	if len(q.data) <= 0 {
		return nil, ErrQueueEmpty
	}
	return q.data,nil

}

// Dequeue's a value from the queue.
func (q *FIFOQueue[T]) Dequeue() (T, error) {
	queue := q.data
	if len(q.data) > 0 {
		card := queue[0]
		q.data = queue[1:]
		return card, nil
	}

	var empty T
	return empty, ErrQueueEmpty
}
func (q *FIFOQueue[T]) Resize(s int) {
	if s > q.Len() {
		q.data = q.data[:s]
	}
}
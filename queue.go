package main

import (
	"github.com/beeker1121/goque"
	"github.com/inloop/go-transport-queue/model"
)

// Queue ...
type Queue struct {
	queue     *goque.PriorityQueue
	transport model.Transport
}

// NewQueue ...
func NewQueue(path string, transport model.Transport) *Queue {
	queue, err := goque.OpenPriorityQueue(path, goque.DESC)

	if err != nil {
		panic(err)
	}

	q := &Queue{queue: queue, transport: transport}
	return q
}

// Push ...
func (q *Queue) Push(priority uint8, message model.TransportMessage) error {
	_, err := q.queue.EnqueueObject(priority, message)
	return err
}

// Pop ...
func (q *Queue) Pop() (*model.TransportMessage, error) {
	item, err := q.queue.Dequeue()
	if err != nil {
		return nil, err
	}
	message, err := q.transport.DecodeMessage(item)
	return &message, err
}

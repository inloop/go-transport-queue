package main

import (
	"time"

	"github.com/inloop/go-transport-queue/model"
)

// MessageBuffer ...
type MessageBuffer struct {
	channel chan model.TransportMessage
	queue   *Queue
	size    int
}

// NewMessageBuffer ...
func NewMessageBuffer(q *Queue, size int) *MessageBuffer {
	c := make(chan model.TransportMessage, size)
	return &MessageBuffer{channel: c, queue: q, size: size}
}

// Start ...
func (b *MessageBuffer) Start(rate time.Duration) {
	go func() {
		throttle := time.Tick(rate)
		for range throttle {
			b.fillBatch()
		}
	}()
}

// GetBatch ...
func (b *MessageBuffer) GetBatch() []model.TransportMessage {
	array := []model.TransportMessage{}
	for {
		select {
		case item := <-b.channel:
			array = append(array, item)
			continue
		default:
			break
		}
		break
	}

	if len(array) == 0 {
		array = append(array, <-b.channel)
	}

	return array
}

func (b *MessageBuffer) fillBatch() {
	for i := 0; i < b.size; i++ {
		item, _ := b.queue.Pop()
		if item != nil {
			b.channel <- *item
		} else {
			break
		}
	}
}

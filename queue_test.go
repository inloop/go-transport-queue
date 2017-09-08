package main

import (
	"testing"

	"github.com/inloop/go-transport-queue/transports"
)

func TestQueueData(t *testing.T) {
	transport := transports.NewLogTransport()
	message := transports.LogTransportMessage{Message: "test"}
	q := NewQueue("queue_test", transport)

	q.Push(0, message)

	item, err := q.Pop()
	if err != nil {
		t.Error(err)
	}

	message2 := (*item).(transports.LogTransportMessage)

	if message.Message != message2.Message {
		t.Errorf("Unexpected message received %s == %s", message.Message, message2.Message)
	}

}

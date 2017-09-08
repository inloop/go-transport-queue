package main

import (
	"testing"

	"github.com/inloop/go-transport-queue/transports"
)

func TestMessageBuffer(t *testing.T) {
	transport := transports.NewLogTransport()
	q := NewQueue("buffer_test", transport)
	b := NewMessageBuffer(q, 2)

	message := transports.LogTransportMessage{Message: "1"}
	message2 := transports.LogTransportMessage{Message: "2"}
	message3 := transports.LogTransportMessage{Message: "3"}

	q.Push(0, message)
	q.Push(0, message2)
	q.Push(1, message3)

	b.fillBatch()

	messages := b.GetBatch()

	if len(messages) != 2 {
		t.Errorf("unexpected messages count %d", len(messages))
	}

	firstMessage := messages[0].(transports.LogTransportMessage)
	if firstMessage.Message != message3.Message {
		t.Errorf("unexpected first message %s", firstMessage.Message)
	}

	secondMessage := messages[1].(transports.LogTransportMessage)
	if secondMessage.Message != message.Message {
		t.Errorf("unexpected second message %s", secondMessage.Message)
	}
}

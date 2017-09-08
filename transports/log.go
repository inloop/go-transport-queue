package transports

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"log"

	"github.com/beeker1121/goque"
	"github.com/inloop/go-transport-queue/model"
)

// NewLogTransport ...
func NewLogTransport() LogTransport {
	gob.Register(LogTransportMessage{})
	return LogTransport{}
}

// LogTransport ...
type LogTransport struct {
}

// DecodeMessages ...
func (t LogTransport) DecodeMessages(d *json.Decoder) ([]model.TransportMessage, error) {
	var message LogTransportMessage
	err := d.Decode(&message)
	if err != nil {
		return []model.TransportMessage{message}, err
	}
	return []model.TransportMessage{message}, nil
}

// SendMessages ...
func (t LogTransport) SendMessages(messages []model.TransportMessage) error {
	for _, message := range messages {
		m := message.(LogTransportMessage)
		log.Println(m.Message)
	}
	return nil
}

// DecodeMessage ...
func (t LogTransport) DecodeMessage(i *goque.PriorityItem) (model.TransportMessage, error) {
	var message LogTransportMessage
	err := i.ToObject(&message)
	return message, err
}

// LogTransportMessage ...
type LogTransportMessage struct {
	Message string `json:"message"`
}

// Validate ...
func (m LogTransportMessage) Validate() error {
	if m.Message == "" {
		return errors.New("message attribute must not be empty")
	}
	return nil
}

package transports

import (
	"encoding/gob"
	"errors"
	"log"

	"github.com/beeker1121/goque"
	"github.com/gin-gonic/gin"
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

// BindResponse ...
func (t LogTransport) BindResponse(c *gin.Context) ([]model.TransportMessage, error) {
	var message LogTransportMessage

	err := c.Bind(&message)
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

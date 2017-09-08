package model

import (
	"encoding/json"

	"github.com/beeker1121/goque"
)

// Transport ...
type Transport interface {
	DecodeMessages(d *json.Decoder) ([]TransportMessage, error)
	SendMessages(messages []TransportMessage) error
	DecodeMessage(i *goque.PriorityItem) (TransportMessage, error)
}

// TransportMessage ...
type TransportMessage interface {
	Validate() error
}

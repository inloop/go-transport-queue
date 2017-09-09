package model

import (
	"github.com/beeker1121/goque"
	"github.com/gin-gonic/gin"
)

// Transport ...
type Transport interface {
	BindResponse(c *gin.Context) ([]TransportMessage, error)
	SendMessages(messages []TransportMessage) error
	DecodeMessage(i *goque.PriorityItem) (TransportMessage, error)
}

// TransportMessage ...
type TransportMessage interface {
	Validate() error
}

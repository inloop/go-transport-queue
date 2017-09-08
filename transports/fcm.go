package transports

import (
	"encoding/gob"
	"encoding/json"
	"fmt"

	"github.com/NaySoftware/go-fcm"
	"github.com/beeker1121/goque"
	"github.com/inloop/go-transport-queue/model"
	uuid "github.com/satori/go.uuid"
)

// "cc9uXlQxRoQ:APA91bFA9ZndOzkXix2GHIjzv1avpCOZSicz9C2_vGkySay_ZqIiFT9LOKPApppWXb4Fg22JeqQRjHyKGHHrIGajR7MThVZYbYhqdn0NxcOsexfhhPPljstliQ-fdGtxWOEzTxTGbolD"

// NewFCMTransport ...
func NewFCMTransport(apiKey string) FCMTransport {
	gob.Register(FCMTransportMessage{})
	return FCMTransport{apiKey: apiKey}
}

// FCMTransport ...
type FCMTransport struct {
	apiKey string
}

// DecodeMessages ...
func (t FCMTransport) DecodeMessages(d *json.Decoder) ([]model.TransportMessage, error) {
	result := []model.TransportMessage{}
	var message FCMTransportMessage
	err := d.Decode(&message)
	if err != nil {
		return result, err
	}

	groupID := uuid.NewV4().String()
	recipients := message.Recipients
	for _, recipient := range recipients {
		message.Recipients = []string{recipient}
		message.GroupID = groupID
		result = append(result, message)
	}

	return result, nil
}

// SendMessages ...
func (t FCMTransport) SendMessages(messages []model.TransportMessage) error {

	messageGroups := map[string]([]FCMTransportMessage){}

	for _, message := range messages {
		m := message.(FCMTransportMessage)
		messageGroups[m.GroupID] = append(messageGroups[m.GroupID], m)
	}
	for _, msgs := range messageGroups {
		firstMessage := msgs[0]
		recipients := []string{}
		for _, msg := range msgs {
			recipients = append(recipients, msg.Recipients[0])
		}
		firstMessage.Recipients = recipients
		if err := t.sendMessage(firstMessage); err != nil {
			return err
		}
	}
	return nil
}

func (t FCMTransport) sendMessage(message FCMTransportMessage) error {
	fmt.Println("Queue: sending fcm", message.Recipients)
	client := fcm.NewFcmClient(t.apiKey)

	client.NewFcmRegIdsMsg(message.Recipients, message.Data)
	client.SetNotificationPayload(&message.Notification)

	_, err := client.Send()
	return err
}

// DecodeMessage ...
func (t FCMTransport) DecodeMessage(i *goque.PriorityItem) (model.TransportMessage, error) {
	var message FCMTransportMessage
	err := i.ToObject(&message)
	return message, err
}

// FCMTransportMessage ...
type FCMTransportMessage struct {
	GroupID      string
	Recipients   []string                `json:"recipients"`
	Notification fcm.NotificationPayload `json:"notification"`
	Data         map[string]string       `json:"data"`
}

// Validate ...
func (m FCMTransportMessage) Validate() error {
	return nil
}

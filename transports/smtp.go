package transports

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/mail"
	"net/url"
	"strconv"

	gomail "gopkg.in/gomail.v2"

	"github.com/beeker1121/goque"
	"github.com/inloop/go-transport-queue/model"
)

// NewSMTPTransport ...
func NewSMTPTransport(urlString, sender string) SMTPTransport {

	gob.Register(SMTPTransportMessage{})
	URL, _ := url.Parse(urlString)

	if URL == nil {
		panic("SMTP url not provided")
	}
	if URL.User == nil {
		panic("user credentials not provided")
	}

	host := URL.Host
	username := URL.User.Username()
	password := ""
	if pass, exists := URL.User.Password(); exists == true {
		password = pass
	}
	port := 25
	if portValue, err := strconv.ParseInt(URL.Port(), 10, 32); err == nil {
		port = int(portValue)
	}

	_sender, err := mail.ParseAddress(sender)
	if err != nil {
		panic(err)
	}

	d := gomail.NewDialer(host, port, username, password)

	return SMTPTransport{dialer: d, sender: *_sender}
}

// SMTPTransport ...
type SMTPTransport struct {
	dialer *gomail.Dialer
	sender mail.Address
}

// DecodeMessages ...
func (t SMTPTransport) DecodeMessages(d *json.Decoder) ([]model.TransportMessage, error) {
	var message SMTPTransportMessage
	err := d.Decode(&message)
	if err != nil {
		return []model.TransportMessage{message}, err
	}
	return []model.TransportMessage{message}, nil
}

// SendMessages ...
func (t SMTPTransport) SendMessages(messages []model.TransportMessage) error {
	for _, message := range messages {
		m := message.(SMTPTransportMessage)
		if err := t.sendMessage(m); err != nil {
			return err
		}
	}
	return nil
}

func (t SMTPTransport) sendMessage(msg SMTPTransportMessage) error {
	fmt.Println("Queue: sending smpt", msg.Recipients)

	m := gomail.NewMessage()
	m.SetHeader("From", t.sender.String())
	m.SetHeader("To", msg.Recipients...)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/plain", msg.Text)
	m.SetBody("text/html", msg.HTML)
	// m.Attach("/home/Alex/lolcat.jpg")
	return t.dialer.DialAndSend(m)
}

// DecodeMessage ...
func (t SMTPTransport) DecodeMessage(i *goque.PriorityItem) (model.TransportMessage, error) {
	var message SMTPTransportMessage
	err := i.ToObject(&message)
	return message, err
}

// SMTPTransportMessage ...
type SMTPTransportMessage struct {
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Text       string   `json:"text"`
	HTML       string   `json:"html"`
}

// Validate ...
func (m SMTPTransportMessage) Validate() error {
	return nil
}

package transports

import (
	"encoding/gob"
	"fmt"
	"net/mail"
	"net/url"
	"strconv"
	"strings"

	gomail "gopkg.in/gomail.v2"

	"github.com/beeker1121/goque"
	"github.com/gin-gonic/gin"
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

	host := strings.Split(URL.Host, ":")[0]
	username := URL.User.Username()
	password := ""
	if pass, exists := URL.User.Password(); exists == true {
		password = pass
	}
	port := 25
	if portValue, err := strconv.ParseInt(URL.Port(), 10, 32); err == nil {
		port = int(portValue)
	}

	d := gomail.NewDialer(host, port, username, password)

	transport := SMTPTransport{dialer: d}
	transport.sender = sender
	return transport
}

// SMTPTransport ...
type SMTPTransport struct {
	dialer *gomail.Dialer
	sender string
}

// BindResponse ...
func (t SMTPTransport) BindResponse(c *gin.Context) ([]model.TransportMessage, error) {
	var message SMTPTransportMessage
	err := c.Bind(&message)
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
		fmt.Println("message sent")
	}
	return nil
}

func (t SMTPTransport) sendMessage(msg SMTPTransportMessage) error {
	fmt.Println("Queue: sending smpt", msg.To)

	sender := msg.From
	if sender == "" {
		sender = t.sender
	}

	address, err := mail.ParseAddress(sender)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", address.Address, address.Name)
	m.SetHeader("To", msg.To...)
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
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
	HTML    string   `json:"html"`
}

// Validate ...
func (m SMTPTransportMessage) Validate() error {
	return nil
}

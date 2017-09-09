package transports

import (
	"encoding/gob"
	"fmt"
	"net/url"
	"strconv"

	gomail "gopkg.in/gomail.v2"

	"github.com/beeker1121/goque"
	"github.com/gin-gonic/gin"
	"github.com/inloop/go-transport-queue/model"
)

// NewSMTPTransport ...
func NewSMTPTransport(urlString string) SMTPTransport {

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

	d := gomail.NewDialer(host, port, username, password)

	return SMTPTransport{dialer: d}
}

// SMTPTransport ...
type SMTPTransport struct {
	dialer *gomail.Dialer
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
	}
	return nil
}

func (t SMTPTransport) sendMessage(msg SMTPTransportMessage) error {
	fmt.Println("Queue: sending smpt", msg.To)

	m := gomail.NewMessage()
	m.SetHeader("From", msg.From)
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

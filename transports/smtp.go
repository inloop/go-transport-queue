package transports

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/mail"
	"net/smtp"
	"net/url"
	"strings"

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
	identity := ""
	auth := smtp.PlainAuth(identity, username, password, host)

	_sender, err := mail.ParseAddress(sender)
	if err != nil {
		panic(err)
	}

	return SMTPTransport{auth: auth, url: *URL, sender: *_sender}
}

// SMTPTransport ...
type SMTPTransport struct {
	auth   smtp.Auth
	sender mail.Address
	url    url.URL
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
		log.Println(m.Message)
		if err := t.sendMessage(m); err != nil {
			return err
		}
	}
	return nil
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	return mime.QEncoding.Encode("utf-8", String)
}

func (t SMTPTransport) sendMessage(msg SMTPTransportMessage) error {
	to := msg.Recipients

	header := make(map[string]string)
	header["From"] = t.sender.String()
	header["To"] = strings.Join(msg.Recipients, ",")
	header["Subject"] = encodeRFC2047(msg.Subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(msg.Message))

	port := t.url.Port()
	if port == "" {
		port = "25"
	}
	return smtp.SendMail(t.url.Hostname()+":"+port, t.auth, t.sender.Address, to, []byte(message))
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
	Message    string   `json:"message"`
}

// Validate ...
func (m SMTPTransportMessage) Validate() error {
	return nil
}

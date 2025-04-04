// internal/smtpclient/smtpclient.go
package smtpclient

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

// Message represents an email message.
type Message struct {
	From    string            // sender email address
	To      []string          // recipient email addresses
	Subject string            // email subject
	Body    string            // plain text body (for simplicity)
	Headers map[string]string // additional headers, e.g. In-Reply-To
}

// NewMessage creates a new Message with initialized headers.
func NewMessage() *Message {
	return &Message{
		Headers: make(map[string]string),
	}
}

// Send constructs the raw message and sends it via the SMTP server.
// For our purposes, we assume Mailpit is running on localhost:1025.
func Send(msg *Message) error {
	if msg.From == "" || len(msg.To) == 0 {
		return errors.New("missing From or To fields")
	}

	var headers []string
	headers = append(headers, fmt.Sprintf("From: %s", msg.From))
	headers = append(headers, fmt.Sprintf("To: %s", strings.Join(msg.To, ", ")))
	headers = append(headers, fmt.Sprintf("Subject: %s", msg.Subject))
	for k, v := range msg.Headers {
		headers = append(headers, fmt.Sprintf("%s: %s", k, v))
	}
	// Basic MIME headers (adjust if you want to send HTML, etc.)
	headers = append(headers, "MIME-Version: 1.0")
	headers = append(headers, "Content-Type: text/plain; charset=\"utf-8\"")

	// Combine headers and body
	raw := []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + msg.Body)

	// Connect to Mailpit’s SMTP server (default localhost:1025)
	addr := "localhost:1025"
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Mail(msg.From); err != nil {
		return err
	}
	for _, addr := range msg.To {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	_, err = bytes.NewBuffer(raw).WriteTo(wc)
	return err
}

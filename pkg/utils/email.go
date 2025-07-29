package utils

import "github.com/go-mail/mail/v2"

type Email struct {
	From     string
	To       string
	Subject  string
	BodyType string
	Body     string
}

func (e *Email) Send() error {
	m := mail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", e.To)
	m.SetHeader("Subject", e.Subject)
	m.SetBody(e.BodyType, e.Body)

	dialer := mail.NewDialer("localhost", 587, "", "")
	err := dialer.DialAndSend(m)
	if err != nil {
		return ErrorHandler(err, "Failed to send reset password")
	}
	return nil
}

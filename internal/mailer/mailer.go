package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func (m *Mailer) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	msg := []byte("From: " + m.From + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" + body)

	return smtp.SendMail(addr, auth, m.From, []string{to}, msg)
}

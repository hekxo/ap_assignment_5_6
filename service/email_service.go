package service

import "gopkg.in/gomail.v2"

type EmailService struct {
	Dialer *gomail.Dialer
	From   string
}

func NewEmailService(host string, port int, username, password, from string) *EmailService {
	dialer := gomail.NewDialer(host, port, username, password)
	return &EmailService{
		Dialer: dialer,
		From:   from,
	}
}

func (s *EmailService) Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.example.com", 587, "user", "password")
	return d.DialAndSend(m)
}

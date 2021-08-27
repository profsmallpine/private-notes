package email

import "net/smtp"

type Service struct {
	email    string
	host     string
	password string
	port     string
}

func NewService(email, host, password, port string) *Service {
	return &Service{
		email:    email,
		host:     host,
		password: password,
		port:     port,
	}
}

func (s *Service) Send(msg, subject string, recipients []string) error {
	rfc822Email := []byte("Subject: " + subject + "\r\n" + "\r\n" + msg + "\r\n") // https://www.rfc-editor.org/rfc/rfc822.html
	auth := smtp.PlainAuth("", s.email, s.password, s.host)
	return smtp.SendMail(s.host+":"+s.port, auth, s.email, recipients, rfc822Email)
}

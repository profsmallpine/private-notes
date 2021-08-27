package domain

type EmailService interface {
	Send(msg, subject string, recipients []string) error
}

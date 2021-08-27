package domain

import "github.com/xy-planning-network/trails/http/session"

type Services struct {
	Auth         AuthService
	Email        EmailService
	SessionStore session.SessionStorer
}

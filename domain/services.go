package domain

import "github.com/xy-planning-network/trails/http/session"

type Services struct {
	Auth         AuthService
	SessionStore session.SessionStorer
}

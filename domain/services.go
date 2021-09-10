package domain

import (
	"github.com/xy-planning-network/trails/http/session"
	"github.com/xy-planning-network/trails/logger"
)

type Services struct {
	Auth         AuthService
	Email        EmailService
	Logger       logger.Logger
	SessionStore session.SessionStorer
}

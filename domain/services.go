package domain

import "github.com/xy-planning-network/trails/logger"

type Services struct {
	Auth      AuthService
	Email     EmailService
	Logger    logger.Logger
	SSE       SSEService
	Websocket WebsocketService
}

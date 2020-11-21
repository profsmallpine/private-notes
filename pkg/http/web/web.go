package web

import (
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/pkg/http/middleware"
	"gorm.io/gorm"
)

type Controller struct {
	*gorm.DB
	*middleware.Manager
	domain.Procedures
	domain.Services
}

func NewController(db *gorm.DB, procedures domain.Procedures, services domain.Services, authKey, encryptionKey []byte) *Controller {
	manager := middleware.NewManager(authKey, encryptionKey)

	return &Controller{
		DB:         db,
		Manager:    manager,
		Procedures: procedures,
		Services:   services,
	}
}

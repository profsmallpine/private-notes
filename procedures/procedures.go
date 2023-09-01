package procedures

import (
	"github.com/profsmallpine/private-notes/domain"
	"gorm.io/gorm"
)

var (
	database = &gorm.DB{}
	services = domain.Services{}
)

func New(allowedEmails []string, db *gorm.DB, s domain.Services) domain.Procedures {
	database = db
	services = s

	return domain.Procedures{
		Meeting: &Meeting{},
		User:    &User{allowedEmails},
	}
}

package migrations

import (
	"gorm.io/gorm"
)

func AddIsAdminToUsers(db *gorm.DB) error {
	return db.Exec("ALTER TABLE users ADD is_admin BOOLEAN DEFAULT false;").Error
}

package postgres

import (
	"github.com/profsmallpine/private-notes/pkg/postgres/migration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(connStr string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := migration.Up(db); err != nil {
		return nil, err
	}

	return db, nil
}

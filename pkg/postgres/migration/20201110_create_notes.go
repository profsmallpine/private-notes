package migration

import (
	"gorm.io/gorm"
)

// Specify table name - default is struct name underscored + pluralized
func (Note) TableName() string {
	return "notes"
}

// Note struct to create notes table
type Note struct {
	gorm.Model
	Content        string
	EncryptionType string
	GroupID        uint
	Group          Group
	Title          string
	UserID         uint
	User           User
}

func CreateNotesTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&Note{})
}

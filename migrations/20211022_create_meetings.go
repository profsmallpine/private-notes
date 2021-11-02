package migrations

import (
	"gorm.io/gorm"
)

// Specify table name - default is struct name underscored + pluralized
func (Meeting) TableName() string {
	return "meetings"
}

// Meeting struct to create meetings table
type Meeting struct {
	gorm.Model
	GroupID uint
	Group   Group
	Status  string
}

func CreateMeetingsTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&Meeting{})
}

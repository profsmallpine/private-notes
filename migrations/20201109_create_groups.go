package migrations

import (
	"gorm.io/gorm"
)

// Specify table name - default is struct name underscored + pluralized
func (Group) TableName() string {
	return "groups"
}

// Group struct to create groups table
type Group struct {
	gorm.Model
	Description string
	Name        string
}

func CreateGroupsTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&Group{})
}

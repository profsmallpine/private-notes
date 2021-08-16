package migrations

import (
	"gorm.io/gorm"
)

// Specify table name - default is struct name underscored + pluralized
func (Comment) TableName() string {
	return "comments"
}

// Comment struct to create comments table
type Comment struct {
	gorm.Model
	Content string
	NoteID  uint
	Note    Note
	UserID  uint
	User    User
}

func CreateCommentsTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&Comment{})
}

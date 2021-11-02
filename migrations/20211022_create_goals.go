package migrations

import (
	"gorm.io/gorm"
)

// Specify table name - default is struct name underscored + pluralized
func (Goal) TableName() string {
	return "goals"
}

// Goal struct to create goals table
type Goal struct {
	gorm.Model
	Completed bool
	Content   string
	MeetingID uint
	Meeting   Meeting
	Mood      string
	UserID    uint
	User      User
}

func CreateGoalsTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&Goal{})
}

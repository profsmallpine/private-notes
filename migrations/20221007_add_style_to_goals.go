package migrations

import (
	"github.com/profsmallpine/private-notes/domain"
	"gorm.io/gorm"
)

// Specify table name - default is struct name underscored + pluralized
func (GoalTable2) TableName() string {
	return "goals"
}

// GoalTable2 struct to update goals table
type GoalTable2 struct {
	gorm.Model
	Completed bool
	Content   string
	MeetingID uint
	Meeting   Meeting
	Mood      string
	Style     domain.GoalStyle `gorm:"type:goal_style"`
	UserID    uint
	User      User
}

func AddStyleToGoals(tx *gorm.DB) error {
	if err := tx.Exec("CREATE TYPE goal_style AS ENUM ('start', 'stop', 'continue');").Error; err != nil {
		return err
	}

	return tx.Migrator().AddColumn(&GoalTable2{}, "style")
}

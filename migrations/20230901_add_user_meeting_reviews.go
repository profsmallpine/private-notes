package migrations

import (
	"gorm.io/gorm"
)

// Specify table name - default is struct name underscored + pluralized
func (UserMeetingReview) TableName() string {
	return "user_meeting_reviews"
}

// UserMeetingReview struct to create user_meeting_reviews table
type UserMeetingReview struct {
	gorm.Model

	MeetingID uint
	Meeting   Meeting
	UserID    uint
	User      User
}

func CreateUserMeetingReviewsTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&UserMeetingReview{})
}

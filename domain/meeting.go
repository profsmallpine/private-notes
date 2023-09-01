package domain

import "gorm.io/gorm"

const MeetingComplete = "complete"

type Meeting struct {
	gorm.Model `json:"-"`

	GroupID uint   `json:"groupID"`
	Status  string `json:"status"`

	Goals []*Goal `json:"goals"`
	Group *Group  `json:"group"`
}

type UserMeetingReview struct {
	gorm.Model `json:"-"`

	MeetingID uint `json:"meetingID"`
	UserID    uint `json:"userID"`
}

type MeetingProcedures interface {
	CopyGoals(*Meeting, *User, []uint) error
	HasPendingReview(*Meeting, *User) (bool, error)
}

func (m *Meeting) CreatedAtHumanized() string {
	return m.CreatedAt.Format("Jan 2, 2006")
}

func (m *Meeting) IsComplete() bool {
	return m.Status == MeetingComplete
}

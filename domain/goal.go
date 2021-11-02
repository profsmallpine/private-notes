package domain

import "gorm.io/gorm"

var GoalMoods = []string{
	"bleak",
	"gloomy",
	"joyful",
	"shitty",
	"hopeful",
	"confident",
	"cheery",
	"expectant",
}

type Goal struct {
	gorm.Model `json:"-"`

	Completed bool   `json:"completed"`
	Content   string `json:"content"`
	MeetingID uint   `json:"meetingID"`
	Mood      string `json:"mood"`
	UserID    uint   `json:"userID"`

	Meeting *Meeting `json:"meeting"`
	User    *User    `json:"user"`
}

package domain

import (
	"gorm.io/gorm"
)

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

	Completed bool      `json:"completed"`
	Content   string    `json:"content"`
	MeetingID uint      `json:"meetingID"`
	Mood      string    `json:"mood"`
	Style     GoalStyle `json:"style"`
	UserID    uint      `json:"userID"`

	Meeting *Meeting `json:"meeting"`
	User    *User    `json:"user"`
}

type GoalStyle string

const (
	GoalStart    GoalStyle = "start"
	GoalStop     GoalStyle = "stop"
	GoalContinue GoalStyle = "continue"
)

var GoalStyles = []GoalStyle{
	GoalStart,
	GoalStop,
	GoalContinue,
}

func (gs GoalStyle) Color() string {
	switch gs {
	case GoalStart:
		return "green"
	case GoalStop:
		return "red"
	case GoalContinue:
		return "purple"
	default:
		return "blue"
	}
}

func (gs GoalStyle) String() string { return string(gs) }

func (gs GoalStyle) Valid() error {
	switch gs {
	case GoalStart, GoalStop, GoalContinue:
		return nil
	default:
		return ErrNotValid
	}
}

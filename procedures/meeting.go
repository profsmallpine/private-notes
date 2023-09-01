package procedures

import (
	"errors"

	"github.com/profsmallpine/private-notes/domain"
	"gorm.io/gorm"
)

type Meeting struct{}

func (m *Meeting) CopyGoals(meeting *domain.Meeting, user *domain.User, goalIDs []uint) error {
	if len(goalIDs) == 0 {
		umr := &domain.UserMeetingReview{
			MeetingID: meeting.ID,
			UserID:    user.ID,
		}
		if err := database.Create(umr).Error; err != nil {
			return err
		}

		return nil
	}

	goals := []*domain.Goal{}
	if err := database.Find(&goals, goalIDs).Error; err != nil {
		return err
	}

	for _, g := range goals {
		copy := &domain.Goal{
			Content:   g.Content,
			MeetingID: meeting.ID,
			Mood:      g.Mood,
			Style:     g.Style,
			UserID:    user.ID,
		}

		if err := database.Create(copy).Error; err != nil {
			return err
		}
	}

	umr := &domain.UserMeetingReview{
		MeetingID: meeting.ID,
		UserID:    user.ID,
	}
	if err := database.Create(umr).Error; err != nil {
		return err
	}

	return nil
}

func (u *Meeting) HasPendingReview(meeting *domain.Meeting, user *domain.User) (bool, error) {
	// If previous meeting cannot be found, then there is no pending review
	lastMeeting := &domain.Meeting{}
	if err := database.Where("id < ?", meeting.ID).Order("id DESC").First(&lastMeeting).Error; err != nil {
		return false, nil
	}

	// If user meeting review record is found, then there is no pending review
	umr := &domain.UserMeetingReview{}
	err := database.Where("user_id = ? AND meeting_id = ?", user.ID, meeting.ID).First(umr).Error
	if err == nil {
		return false, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}

	return false, err
}

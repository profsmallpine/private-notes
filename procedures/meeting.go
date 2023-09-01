package procedures

import (
	"github.com/profsmallpine/private-notes/domain"
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
	return true, nil
}

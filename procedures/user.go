package procedures

import (
	"errors"

	"github.com/profsmallpine/private-notes/domain"
	"github.com/xy-planning-network/trails/logger"
	"gorm.io/gorm"
)

type User struct {
	allowedEmails []string
}

func (u *User) FetchWithGroups(id uint, user *domain.User) error {
	return database.Where("id = ?", id).Preload("Groups").First(user).Error
}

func (u *User) HandleCallback(data *domain.AuthData) (*domain.User, error) {
	if data.Email == "" || !u.isAllowedEmail(data.Email) {
		return nil, domain.ErrNotValid
	}

	user := &domain.User{}
	if err := database.Where("email = ?", data.Email).First(user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		user.Email = data.Email
		user.FirstName = data.FirstName
		user.LastName = data.LastName
		user.PictureURL = data.PictureURL

		if err := database.Create(user).Error; err != nil {
			services.Logger.Error(err.Error(), &logger.LogContext{Error: err})
			return nil, err
		}
	}

	return user, nil
}

func (u *User) isAllowedEmail(email string) bool {
	for _, allowed := range u.allowedEmails {
		if email == allowed {
			return true
		}
	}
	return false
}

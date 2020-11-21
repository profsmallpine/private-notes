package domain

import (
	"strings"

	"github.com/profsmallpine/private-notes/pkg/http/routes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`

	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	PictureURL string `json:"pictureURL"`

	Groups []*Group `gorm:"many2many:user_groups;"`
}

type UserProcedures interface {
	HandleCallback(data *AuthData) (*User, error)
}

func (u *User) CanAccessGroup(groupID uint) bool {
	for _, group := range u.Groups {
		if group.ID == groupID {
			return true
		}
	}

	return false
}

func (*User) HomePath() string {
	return routes.GetGroupsURL
}

func (u *User) Initials() string {
	return strings.TrimSpace(string(u.FirstName[0]) + string(u.LastName[0]))
}

func (u *User) FullName() string {
	return strings.TrimSpace(
		strings.Join(
			[]string{u.FirstName, u.LastName}, " ",
		),
	)
}

func NewUserFromID(id uint) *User {
	u := &User{}
	u.ID = id
	return u
}

func (*User) TableName() string {
	return "users"
}

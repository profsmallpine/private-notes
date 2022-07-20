package domain

import (
	"strings"

	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/middleware"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`

	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	IsAdmin    bool   `json:"isAdmin"`
	LastName   string `json:"lastName"`
	PictureURL string `json:"pictureURL"`

	Groups []*Group `gorm:"many2many:user_groups;"`
}

type UserProcedures interface {
	CanAccessGroup(user *User, groupID string) (*Group, error)
	HandleCallback(data *AuthData) (*User, error)
	GetByID(id uint) (middleware.User, error)
}

func (u *User) CanAccessGroup(groupID uint) bool {
	for _, group := range u.Groups {
		if group.ID == groupID {
			return true
		}
	}

	return false
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetID() uint {
	return u.ID
}

func (*User) HasAccess() bool {
	return true
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

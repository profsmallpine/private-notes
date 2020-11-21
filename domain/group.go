package domain

import "gorm.io/gorm"

type Group struct {
	gorm.Model `json:"-"`

	Description string `json:"description"`
	Name        string `json:"name"`

	Users []*User `gorm:"many2many:user_groups;"`
}

func (g *Group) CreatedAtHumanized() string {
	return g.CreatedAt.Format("Jan 2, 2006")
}

func (*Group) TableName() string {
	return "groups"
}

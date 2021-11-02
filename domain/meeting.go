package domain

import "gorm.io/gorm"

type Meeting struct {
	gorm.Model `json:"-"`

	GroupID uint   `json:"groupID"`
	Status  string `json:"status"`

	Goals []*Goal `json:"goals"`
	Group *Group  `json:"group"`
}

func (m *Meeting) CreatedAtHumanized() string {
	return m.CreatedAt.Format("Jan 2, 2006")
}

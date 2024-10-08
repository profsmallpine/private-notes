package domain

import (
	"html/template"

	"gorm.io/gorm"
)

type Note struct {
	gorm.Model `json:"-"`

	Content        string `json:"content"`
	EncryptionType string `json:"encryptionType"`
	GroupID        uint   `json:"groupID"`
	Title          string `json:"title"`
	UserID         uint   `json:"userID"`

	Author   *User      `json:"author" gorm:"foreignKey:UserID"`
	Comments []*Comment `json:"comments"`
	Group    *Group     `json:"group"`
}

func (n *Note) CreatedAtHumanized() string {
	return n.CreatedAt.Format("Jan 2, 2006")
}

func (n *Note) HTMLContent() template.HTML {
	return template.HTML(n.Content)
}

func (*Note) TableName() string {
	return "notes"
}

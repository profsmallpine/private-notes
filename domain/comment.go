package domain

import "gorm.io/gorm"

type Comment struct {
	gorm.Model `json:"-"`

	Content string `json:"content"`
	NoteID  uint   `json:"noteID"`
	UserID  uint   `json:"userID"`

	Author *User `json:"author" gorm:"foreignKey:UserID"`
	Note   *Note `json:"note"`
}

func (c *Comment) CreatedAtHumanized() string {
	return c.CreatedAt.Format("Jan 2, 2006")
}

func (*Comment) TableName() string {
	return "comments"
}

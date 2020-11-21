package domain

import "gorm.io/gorm"

type Note struct {
	gorm.Model `json:"-"`

	Content        string `json:"content"`
	EncryptionType string `json:"encryptionType"`
	GroupID        uint   `json:"groupID"`
	Title          string `json:"title"`
	UserID         uint   `json:"userID"`

	Author   *User      `json:"author" gorm:"foreignKey:UserID"`
	Comments []*Comment `json:"comments"`
	// Group via group_id
}

func (n *Note) CreatedAtHumanized() string {
	return n.CreatedAt.Format("Jan 2, 2006")
}

func (*Note) TableName() string {
	return "notes"
}

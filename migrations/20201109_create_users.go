package migrations

import (
	"gorm.io/gorm"
)

// Specify table name - default is struct name underscored + pluralized
func (User) TableName() string {
	return "users"
}

// User struct to create users table
type User struct {
	gorm.Model
	Email      string `sql:"index"`
	FirstName  string
	LastName   string
	PictureURL string
	Groups     []Group `gorm:"many2many:user_groups;"`
}

func CreateUsersTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&User{})
}

package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID   int64  `gorm:"not null"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Email    string `gorm:"unique" json:"email"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

func (User) TableName() string {
	return "users"
}

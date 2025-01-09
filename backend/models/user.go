package models

import (
	"time"
)

type User struct {
	ID           int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string     `json:"username" gorm:"type:varchar(50);unique;not null"`
	Email        string     `json:"email" gorm:"type:varchar(100);unique;not null"`
	PasswordHash string     `json:"-" gorm:"type:varchar(255);not null;column:password_hash"`
	AvatarURL    string     `json:"avatar_url" gorm:"type:varchar(255);column:avatar_url"`
	Bio          string     `json:"bio" gorm:"type:text"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Status       int8       `json:"status" gorm:"type:tinyint;default:1;comment:'1: active, 0: inactive'"`
	LastLoginAt  *time.Time `json:"last_login_at" gorm:"column:last_login_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

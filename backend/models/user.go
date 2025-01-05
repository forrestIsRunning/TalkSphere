package models

import (
	"time"
)

type User struct {
	ID           int64      `gorm:"primaryKey;autoIncrement"`
	Username     string     `gorm:"type:varchar(50);unique;not null"`
	Email        string     `gorm:"type:varchar(100);unique;not null"`
	PasswordHash string     `gorm:"type:varchar(255);not null;column:password_hash"`
	AvatarURL    string     `gorm:"type:varchar(255);column:avatar_url"`
	Bio          string     `gorm:"type:text"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Status       int8       `gorm:"type:tinyint;default:1;comment:'1: active, 0: inactive'"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

package models

import "time"

// Favorite 收藏模型
type Favorite struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"index" json:"user_id"`
	PostID    int64     `gorm:"index" json:"post_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// 关联字段
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post *Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}

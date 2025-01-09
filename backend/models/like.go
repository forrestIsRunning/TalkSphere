package models

import "time"

// Like 点赞模型
type Like struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"index" json:"user_id"`
	TargetID   int64     `gorm:"index;comment:'点赞目标ID'" json:"target_id"`
	TargetType int8      `gorm:"type:tinyint;comment:'1: post, 2: comment'" json:"target_type"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// 关联字段
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}

// 点赞类型常量
const (
	LikeTargetPost    int8 = 1
	LikeTargetComment int8 = 2
)

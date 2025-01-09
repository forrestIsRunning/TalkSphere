package models

import "time"

// Comment 评论模型
type Comment struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	PostID     int64     `json:"post_id" gorm:"not null"`
	UserID     int64     `json:"user_id" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	ParentID   *int64    `json:"parent_id" gorm:"default:null"`
	RootID     *int64    `json:"root_id" gorm:"default:null"`
	LikeCount  int       `json:"like_count" gorm:"default:0"`
	ReplyCount int       `json:"reply_count" gorm:"default:0"`
	Score      int       `json:"score" gorm:"default:0"`
	Status     int8      `json:"status" gorm:"default:1"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// 非数据库字段
	Children []Comment `gorm:"-" json:"children,omitempty"`
	User     *User     `gorm:"-" json:"user,omitempty"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}

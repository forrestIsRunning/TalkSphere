package models

import "time"

type Post struct {
	ID            int64       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title         string      `json:"title" gorm:"type:varchar(255);not null"`
	Content       string      `json:"content" gorm:"type:text;not null"`
	BoardID       *int64      `json:"board_id" gorm:"column:board_id"`
	AuthorID      *int64      `json:"author_id" gorm:"column:author_id"`
	ViewCount     int         `json:"view_count" gorm:"default:0;comment:'观看次数'"`
	LikeCount     int         `json:"like_count" gorm:"default:0;comment:'点赞数'"`
	FavoriteCount int         `json:"favorite_count" gorm:"default:0;comment:'收藏数'"`
	CommentCount  int         `json:"comment_count" gorm:"default:0;comment:'评论数'"`
	CreatedAt     time.Time   `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time   `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Status        int8        `json:"status" gorm:"type:tinyint;default:1;comment:'1: published, 0: draft, -1: deleted'"`
	Tags          []Tag       `json:"tags" gorm:"many2many:post_tags;"`
	Images        []PostImage `json:"images" gorm:"foreignKey:PostID"`
	Author        *User       `json:"author" gorm:"foreignKey:AuthorID"`
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}

// PostImage 帖子图片模型
type PostImage struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	PostID    int64     `gorm:"index;not null"`
	UserID    int64     `gorm:"index;not null"`
	ImageURL  string    `gorm:"type:varchar(255);not null"`
	Status    int8      `gorm:"default:1"`
	SortOrder int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (PostImage) TableName() string {
	return "post_images"
}

// Tag 标签模型
type Tag struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Posts     []Post    `gorm:"many2many:post_tags;"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

type PostTag struct {
	PostID int64 `gorm:"primaryKey;column:post_id"`
	TagID  int64 `gorm:"primaryKey;column:tag_id"`
}

func (PostTag) TableName() string {
	return "post_tags"
}

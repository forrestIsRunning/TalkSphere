package models

import "time"

type Board struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name;type:varchar(100);not null"`
	Description string    `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
	Status      int8      `gorm:"column:status;type:tinyint;default:1;comment:'1: active, 0: inactive'"`
	SortOrder   int       `gorm:"column:sort_order;default:0"`
	CreatorID   int64     `gorm:"column:creator_id"`
}

// TableName 指定表名
func (Board) TableName() string {
	return "boards"
}

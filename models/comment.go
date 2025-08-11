package models

import "gorm.io/gorm"

//评论表
type Comment struct {
	gorm.Model
	Content string `json:"content" gorm:"type:text;not null"`
	UserID  uint   `gorm:"not null"`          // 外键
	User    User   `gorm:"foreignKey:UserID"` // 属于关系
	PostID  uint   `gorm:"not null"`          // 外键
	Post    Post   `gorm:"foreignKey:PostID"` // 属于关系
}

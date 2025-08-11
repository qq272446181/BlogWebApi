package models

import "gorm.io/gorm"

//博客文章
type Post struct {
	gorm.Model
	Title    string    `json:"title" gorm:"size:255;not null"`
	Content  string    `json:"content" gorm:"type:text;not null"`
	AuthorID uint      `gorm:"not null"`            // 外键
	Author   User      `gorm:"foreignKey:AuthorID"` // 属于关系
	Comments []Comment `gorm:"foreignKey:PostID"`
}

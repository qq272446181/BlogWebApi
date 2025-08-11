package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户表
type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null"`
	Password string `gorm:"size:100;not null"`
	Email    string `gorm:"size:255;uniqueIndex;not null"`
}

func (u *User) HashPassword(inputpwd string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(inputpwd), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}
func (u *User) CheckPassword(inputpwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputpwd))
	if err != nil {
		return err
	}
	return nil
}

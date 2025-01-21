package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id         int64     `gorm:"primaryKey" json:"id"`
	FirstName  string    `gorm:"size:255;not null" json:"firstName"`
	LastName   string    `gorm:"size:255;not null" json:"lastName"`
	Username   string    `gorm:"unique;size:255;not null" json:"username"`
	ChatId     int64     `gorm:"not null" json:"chatId"`
	IsAdmin    bool      `gorm:"default:false" json:"isAdmin"`
	IsApproved bool      `gorm:"default:false" json:"isApproved"`
	IsDeclined bool      `gorm:"default:false" json:"isDeclined"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	return
}

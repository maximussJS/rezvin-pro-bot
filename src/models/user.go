package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id         int64     `gorm:"primaryKey" json:"id"`
	FirstName  string    `gorm:"size:100;not null" json:"firstName"`
	LastName   string    `gorm:"size:100" json:"lastName"`
	Username   string    `gorm:"size:100" json:"username"`
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

func (u *User) GetPrivateName() string {
	text := fmt.Sprintf("%s %s", u.FirstName, u.LastName)

	if u.Username != "" {
		text += fmt.Sprintf(" @%s", u.Username)
	}

	return text
}

func (u *User) GetPublicName() string {
	text := fmt.Sprintf("%s %s", u.FirstName, u.LastName)

	if u.Username != "" {
		text += fmt.Sprintf(" @%s", u.Username)
	}

	return text
}

func (u *User) IsNotApproved() bool {
	return !u.IsApproved
}

func (u *User) IsNotAdmin() bool {
	return !u.IsAdmin
}

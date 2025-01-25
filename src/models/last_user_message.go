package models

import (
	"gorm.io/gorm"
	"time"
)

type LastUserMessage struct {
	MessageId int       `gorm:"index:idx_last_user_message,unique" json:"messageId"`
	ChatId    int64     `gorm:"primaryKey;autoIncrement=false" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt  time.Time `json:"updateAt"`
}

func (u *LastUserMessage) TableName() string {
	return "last_user_messages"
}

func (u *LastUserMessage) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	return
}

func (u *LastUserMessage) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateAt = time.Now()
	return
}

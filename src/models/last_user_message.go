package models

import (
	"fmt"
	"gorm.io/gorm"
	"rezvin-pro-bot/src/config"
	"time"
)

type LastUserMessage struct {
	MessageId int       `gorm:"index:idx_last_user_message,unique;not null" json:"messageId"`
	ChatId    int64     `gorm:"primaryKey;autoIncrement=false" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt  time.Time `json:"updateAt"`
}

func (u *LastUserMessage) TableName() string {
	schema := config.GetPostgresSchema()

	return fmt.Sprintf("%s.last_user_messages", schema)
}

func (u *LastUserMessage) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	return
}

func (u *LastUserMessage) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdateAt = time.Now()
	return
}

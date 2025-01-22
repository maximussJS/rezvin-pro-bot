package models

import (
	"gorm.io/gorm"
	"time"
)

type UserProgram struct {
	UserId    int64     `gorm:"primaryKey" json:"userId"`
	ProgramId uint      `gorm:"primaryKey" json:"programId"`
	User      User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Program   Program   `gorm:"foreignKey:ProgramId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"program"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u *UserProgram) TableName() string {
	return "user_programs"
}

func (u *UserProgram) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	return
}

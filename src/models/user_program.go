package models

import (
	"fmt"
	"gorm.io/gorm"
	"rezvin-pro-bot/src/config"
	"time"
)

type UserProgram struct {
	Id        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"index:idx_user_program,unique;not null" json:"userId"`
	ProgramId uint      `gorm:"index:idx_user_program,unique;not null" json:"programId"`
	User      User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Program   Program   `gorm:"foreignKey:ProgramId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"program"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u *UserProgram) Name() string {
	return u.Program.Name
}

func (u *UserProgram) TableName() string {
	schema := config.GetPostgresSchema()
	return fmt.Sprintf("%s.user_programs", schema)
}

func (u *UserProgram) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	return
}

package models

import (
	"gorm.io/gorm"
	"time"
)

type UserResult struct {
	Id            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserProgramId uint      `gorm:"index:idx_record,unique;not null" json:"userProgramId"`
	ExerciseId    uint      `gorm:"index:idx_record,unique;not null" json:"exerciseId"`
	Reps          uint      `gorm:"index:idx_record,unique;not null" json:"reps"`
	Weight        int       `gorm:"not null" json:"weight"`
	Exercise      Exercise  `gorm:"foreignKey:ExerciseId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"exercise"`
	LoggedAt      time.Time `json:"loggedAt"`
}

func (u *UserResult) Name() string {
	return u.Exercise.Name
}

func (u *UserResult) TableName() string {
	return "user_exercise_records"
}

func (u *UserResult) BeforeCreate(tx *gorm.DB) (err error) {
	u.LoggedAt = time.Now()
	return
}

func (u *UserResult) BeforeUpdate(tx *gorm.DB) (err error) {
	u.LoggedAt = time.Now()
	return
}

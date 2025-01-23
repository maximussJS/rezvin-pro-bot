package models

import (
	"gorm.io/gorm"
	"time"
)

type UserExerciseRecord struct {
	Id            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserProgramId uint      `gorm:"index:idx_record,unique" json:"userProgramId"`
	ExerciseId    uint      `gorm:"index:idx_record,unique" json:"exerciseId"`
	Reps          uint      `gorm:"index:idx_record,unique" json:"reps"`
	Weight        int       `gorm:"not null" json:"weight"`
	Exercise      Exercise  `gorm:"foreignKey:ExerciseId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"exercise"`
	LoggedAt      time.Time `json:"loggedAt"`
}

func (u *UserExerciseRecord) Name() string {
	return u.Exercise.Name
}

func (u *UserExerciseRecord) TableName() string {
	return "user_exercise_records"
}

func (u *UserExerciseRecord) BeforeCreate(tx *gorm.DB) (err error) {
	u.LoggedAt = time.Now()
	return
}

func (u *UserExerciseRecord) BeforeUpdate(tx *gorm.DB) (err error) {
	u.LoggedAt = time.Now()
	return
}

package models

import (
	"gorm.io/gorm"
	"time"
)

type UserMeasure struct {
	Id        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"index:idx_user_program,unique;not null" json:"userId"`
	MeasureId uint      `gorm:"index:idx_measure;not null" json:"measureId"`
	User      User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Measure   Measure   `gorm:"foreignKey:MeasureId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"measure"`
	Value     float64   `gorm:"not null" json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *UserMeasure) Name() string {
	return u.Measure.Name
}

func (u *UserMeasure) TableName() string {
	return "user_measures"
}

func (u *UserMeasure) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return
}

func (u *UserMeasure) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}

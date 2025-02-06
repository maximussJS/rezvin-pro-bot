package models

import (
	"fmt"
	"gorm.io/gorm"
	"rezvin-pro-bot/src/globals"
	"time"
)

type UserMeasure struct {
	Id        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"index:idx_user_measure_user_id;not null" json:"userId"`
	MeasureId uint      `gorm:"index:idx_user_measure_measure_id;not null" json:"measureId"`
	Value     float64   `gorm:"not null" json:"value"`
	User      User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Measure   Measure   `gorm:"foreignKey:MeasureId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"measure"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u *UserMeasure) Name() string {
	return u.Measure.Name
}

func (u *UserMeasure) Units() string {
	return u.Measure.Units
}

func (u *UserMeasure) TableName() string {
	schema := globals.GetPostgresSchema()
	return fmt.Sprintf("%s.user_measures", schema)
}

func (u *UserMeasure) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	return
}

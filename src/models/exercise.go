package models

import (
	"gorm.io/gorm"
	"time"
)

type Exercise struct {
	Id        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"index:idx_exercise,unique;size:100;not null" json:"name"`
	ProgramId uint      `gorm:"not null;index:idx_exercise,unique;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"programId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (p *Exercise) TableName() string {
	return "exercises"
}

func (p *Exercise) BeforeCreate(tx *gorm.DB) (err error) {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return
}

func (p *Exercise) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return
}

package models

import (
	"fmt"
	"gorm.io/gorm"
	"rezvin-pro-bot/src/config"
	"time"
)

type Program struct {
	Id        uint       `gorm:"primaryKey;autoIncrement" json:"id" `
	Name      string     `gorm:"size:100;not null;unique" json:"name"`
	Exercises []Exercise `gorm:"foreignKey:ProgramId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"exercises"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (c *Program) TableName() string {
	schema := config.GetPostgresSchema()
	return fmt.Sprintf("%s.programs", schema)
}

func (c *Program) BeforeCreate(tx *gorm.DB) (err error) {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *Program) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}

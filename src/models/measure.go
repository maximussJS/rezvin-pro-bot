package models

import (
	"gorm.io/gorm"
	"time"
)

type Measure struct {
	Id        uint      `gorm:"primaryKey;autoIncrement" json:"id" `
	Name      string    `gorm:"size:100;not null;unique" json:"name"`
	Units     string    `gorm:"size:100;not null" json:"units"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c *Measure) TableName() string {
	return "measures"
}

func (c *Measure) BeforeCreate(tx *gorm.DB) (err error) {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *Measure) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}

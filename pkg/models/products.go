package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name     string        `gorm:"not null"`
	Price    float32       `gorm:"not null"`
	Duration time.Duration `gorm:"not null"`
}

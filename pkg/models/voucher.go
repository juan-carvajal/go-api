package models

import (
	"time"

	"gorm.io/gorm"
)

type Voucher struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ValidThru time.Time
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID             uint           `gorm:"primarykey"`
	CreatedAt      time.Time      `gorm:"default:NOW()"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	UserID         uint           `gorm:"not null;uniqueIndex:idx_unique_subscription"`
	User           User           `gorm:"foreignKey:UserID;references:ID" json:"-"`
	VoucherID      *string
	Voucher        *Voucher `gorm:"foreignKey:VoucherID;references:ID" json:"-"`
	StartDate      time.Time
	EndDate        time.Time
	ProductID      uint    `gorm:"not null;uniqueIndex:idx_unique_subscription"`
	Product        Product `gorm:"foreignKey:ProductID;references:ID" json:"-"`
	Price          float32 `gorm:"not null"`
	IsOnTrial      bool
	TrialEnabled   bool      `json:"-"`
	Paused         bool      `gorm:"default:FALSE"`
	LastPausedDate time.Time `gorm:"not null" json:"-"`
}

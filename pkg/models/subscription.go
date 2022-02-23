package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	UserID       uint `gorm:"not null;uniqueIndex:idx_unique_subscription"`
	User         User `gorm:"foreignKey:UserID;references:ID"`
	VoucherID    string
	Voucher      Voucher `gorm:"foreignKey:VoucherID;references:ID"`
	StartDate    time.Time
	EndDate      time.Time
	ProductID    uint    `gorm:"not null;uniqueIndex:idx_unique_subscription"`
	Product      Product `gorm:"foreignKey:ProductID;references:ID"`
	IsOnTrial    bool
	TrialEnabled bool
}

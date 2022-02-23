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
	ValidThru time.Time      `gorm:"not null"`
	Type      VoucherType    `gorm:"not null"`
	Discount  float32        `gorm:"not null"`
}

type VoucherType string

const (
	Fixed      VoucherType = "fixed"
	Percentage VoucherType = "percentage"
)

type VoucherRedeem struct {
	ID         uint      `gorm:"primarykey"`
	RedeemedAt time.Time `gorm:"not null"`
	UserID     uint
	VoucherID  string
	User       User    `gorm:"foreignKey:UserID;references:ID"`
	Voucher    Voucher `gorm:"foreignKey:VoucherID;references:ID"`
}

package models

import (
	"encoding/json"
	"errors"
	"math"
	"time"

	"gorm.io/gorm"
)

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type Product struct {
	ID        uint       `gorm:"primarykey"`
	CreatedAt *time.Time `gorm:"default:NOW()"`
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"not null"`
	Price     float32        `gorm:"not null"`
	Duration  Duration       `gorm:"not null"`
}

func (p *Product) DiscountedPrice(voucher Voucher) float32 {
	if voucher.Type == Fixed {
		return float32(math.Max(float64(p.Price)-float64(voucher.Discount), 0))
	} else {
		return p.Price * (1 - voucher.Discount)
	}
}

package voucher

import (
	"github.com/juan-carvajal/go-api/pkg/models"
	"gorm.io/gorm"
)

type VoucherRepo interface {
	GetVoucherByID(id string) (*models.Voucher, error)
	GetVoucherRedeem(userId uint, voucherId string) (*models.VoucherRedeem, error)
}

type DefaultVoucherRepo struct {
	db *gorm.DB
}

func NewDefaultVoucherRepo(db *gorm.DB) VoucherRepo {
	return &DefaultVoucherRepo{db}
}

func (r *DefaultVoucherRepo) GetVoucherByID(id string) (*models.Voucher, error) {
	voucher := &models.Voucher{}

	tx := r.db.First(&voucher, "id = ?", id)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return voucher, nil
}

func (r *DefaultVoucherRepo) GetVoucherRedeem(userId uint, voucherId string) (*models.VoucherRedeem, error) {
	voucherRedeem := &models.VoucherRedeem{}

	tx := r.db.First(&voucherRedeem, "user_id = ? AND voucher_id >= ?", userId, voucherId)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return voucherRedeem, nil
}

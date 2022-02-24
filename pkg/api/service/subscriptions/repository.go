package subscriptions

import (
	"fmt"
	"time"

	"github.com/juan-carvajal/go-api/pkg/models"
	"gorm.io/gorm"
)

type SubscriptionRepo interface {
	GetUserSubscriptions(userId uint) (*[]models.Subscription, error)
	CancelSubscription(id uint) error
	GetSubscriptionById(id uint) (*models.Subscription, error)
	PauseSubscription(id uint) error
	UnpauseSubscription(id uint) error
	CreateSubscription(params CreateSubscriptionParams) error
}

type DefaultSubscriptionRepo struct {
	db *gorm.DB
}

type CreateSubscriptionParams struct {
	UserID  uint
	Voucher *models.Voucher
	Product models.Product
}

func (r *DefaultSubscriptionRepo) GetUserSubscriptions(userId uint) (*[]models.Subscription, error) {
	var subs []models.Subscription

	tx := r.db.Find(&subs, "user_id = ?", userId)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &subs, nil
}

func (r *DefaultSubscriptionRepo) CancelSubscription(id uint) error {

	tx := r.db.Model(&models.Subscription{}).Where("id = ?", id).Update("deleted_at", time.Now())

	return tx.Error
}

func (r *DefaultSubscriptionRepo) GetSubscriptionById(id uint) (*models.Subscription, error) {
	sub := &models.Subscription{}

	tx := r.db.First(&sub, id)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return sub, nil
}

func (r *DefaultSubscriptionRepo) PauseSubscription(id uint) error {
	sub, err := r.GetSubscriptionById(id)
	if err != nil {
		return err
	}

	if sub.IsOnTrial || sub.Paused {
		return nil
	}

	tx := r.db.Model(&models.Subscription{}).Where("id = ?", id).Updates(map[string]interface{}{"paused": true, "last_paused_date": time.Now()})

	return tx.Error
}

func (r *DefaultSubscriptionRepo) UnpauseSubscription(id uint) error {
	sub, err := r.GetSubscriptionById(id)
	if err != nil {
		return err
	}

	if sub.IsOnTrial || !sub.Paused {
		return nil
	}

	newDuration := sub.LastPausedDate.Sub(sub.StartDate)

	tx := r.db.Model(&models.Subscription{}).Where("id = ?", id).Updates(map[string]interface{}{"paused": false, "end_date": sub.EndDate.Add(newDuration)})

	return tx.Error
}

func (r *DefaultSubscriptionRepo) CreateSubscription(params CreateSubscriptionParams) error {
	startDate := time.Now().Add(time.Hour * 24 * 30)

	sub := models.Subscription{
		CreatedAt: time.Now(),
		ProductID: params.Product.ID,
		StartDate: startDate,
		EndDate:   startDate.Add(time.Duration(params.Product.Duration)),
		IsOnTrial: true,
		Price:     params.Product.Price,
		UserID:    params.UserID,
	}

	fmt.Printf("%+v", params)

	if params.Voucher != nil {
		sub.Price = params.Product.DiscountedPrice(*params.Voucher)
		sub.VoucherID = &params.Voucher.ID
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&sub).Error; err != nil {
			return err
		}

		if params.Voucher != nil {
			if err := tx.Create(&models.VoucherRedeem{
				RedeemedAt: time.Now(),
				UserID:     params.UserID,
				VoucherID:  params.Voucher.ID,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func NewDefaultSubscriptionsRepo(db *gorm.DB) SubscriptionRepo {
	return &DefaultSubscriptionRepo{
		db: db,
	}
}

package migrations

import (
	"fmt"
	"time"

	"github.com/juan-carvajal/go-api/pkg/models"
	"gorm.io/gorm"
)

var products = []models.Product{
	{Name: "Product 1", Duration: models.Duration(time.Hour * 24 * 60), Price: 5000},
	{Name: "Product 2", Duration: models.Duration(time.Hour * 24 * 30), Price: 10000},
	{Name: "Product 3", Duration: models.Duration(time.Hour * 24 * 120), Price: 20000},
	{Name: "Product 4", Duration: models.Duration(time.Hour * 24 * 180), Price: 30000},
	{Name: "Product 5", Duration: models.Duration(time.Hour * 24 * 5), Price: 2000},
	{Name: "Product 6", Duration: models.Duration(time.Hour * 24 * 10), Price: 4000},
	{Name: "Product 7", Duration: models.Duration(time.Hour * 24 * 15), Price: 1000},
}

var vouchers = []models.Voucher{
	{ValidThru: time.Now().Add(time.Hour * 24 * 365), Type: models.Fixed, Discount: 500, ID: "test-fixed"},
}

func AutoMigrateAndSeed(db *gorm.DB) error {

	err := db.Debug().AutoMigrate(models.Product{}, models.Voucher{}, models.User{}, models.VoucherRedeem{}, models.Subscription{})
	if err != nil {
		return err
	}

	fmt.Println("migrations completed")

	fmt.Println("seeding products")

	db.Create(&products)

	fmt.Println("seeding vouchers")

	db.Create(&vouchers)

	return nil
}

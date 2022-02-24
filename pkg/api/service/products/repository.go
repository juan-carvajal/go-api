package products

import (
	"time"

	"github.com/juan-carvajal/go-api/pkg/models"
	"github.com/juan-carvajal/go-api/pkg/models/shared"
	"gorm.io/gorm"
)

type ProductQueryParams struct {
	shared.SearchParams
	Voucher *models.Voucher
}

type ProductRepo interface {
	GetAllProducts(params ProductQueryParams) (*[]models.Product, error)
	GetProductByID(id int) (*models.Product, error)
}

type DefaultProductRepo struct {
	db *gorm.DB
}

func (r *DefaultProductRepo) GetAllProducts(params ProductQueryParams) (*[]models.Product, error) {
	var products []models.Product

	query := *r.db

	if params.Search != "" {
		query = *query.Where("name ilike ?", params.Search)

		if query.Error != nil {
			return nil, query.Error
		}
	}

	query = *query.Offset(params.Offset).Limit(params.PageSize)
	if query.Error != nil {
		return nil, query.Error
	}

	query = *query.Find(&products)
	if query.Error != nil {
		return nil, query.Error
	}

	if params.Voucher == nil {
		return &products, nil
	}

	if params.Voucher.ValidThru.Before(time.Now()){
		return &products, nil
	}

	for i := range products {
		products[i].Price = products[i].DiscountedPrice(*params.Voucher)
	}

	// if params.Voucher.Type == models.Fixed {
	// 	for i := range products {
	// 		products[i].Price = float32(math.Max(float64(products[i].Price)-float64(params.Voucher.Discount), 0))
	// 	}
	// } else {
	// 	for i := range products {
	// 		products[i].Price = products[i].Price * (1 - params.Voucher.Discount)
	// 	}
	// }

	return &products, nil
}

func (r *DefaultProductRepo) GetProductByID(id int) (*models.Product, error) {
	product := &models.Product{}

	tx := r.db.First(&product, id)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return product, nil
}

func NewDefaultProductRepo(db *gorm.DB) ProductRepo {
	return &DefaultProductRepo{db}
}

func NameSearch(search string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name IN (?)", search)
	}
}

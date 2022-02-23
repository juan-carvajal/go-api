package products

import (
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
		query = *query.Where("name ilike %?%", params.Search)

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

	return &products, nil
}

func (r *DefaultProductRepo) GetProductByID(id int) (*models.Product, error) {
	return nil, nil
}

func NewDefaultProductRepo(db *gorm.DB) ProductRepo {
	return &DefaultProductRepo{db}
}

func NameSearch(search string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name IN (?)", search)
	}
}

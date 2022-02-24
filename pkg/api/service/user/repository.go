package user

import (
	"github.com/juan-carvajal/go-api/pkg/models"
	"gorm.io/gorm"
)

type UserRepo interface {
	GetUserByUsername(username string) (*models.User, error)
	RegisterUser(user models.User) error
}

type DefaultUserRepo struct {
	db *gorm.DB
}

func (r *DefaultUserRepo) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}

	tx := r.db.First(&user, "username = ?", username)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func (r *DefaultUserRepo) RegisterUser(user models.User) error {
	tx := r.db.Select("Username", "PasswordHash").Create(&user)

	return tx.Error
}

func NewDefaultUserRepo(db *gorm.DB) UserRepo {
	return &DefaultUserRepo{db: db}
}

package repository

import (
	"github.com/kev1nandreas/go-rest-api-template/src/modules/user/entity"
	"gorm.io/gorm"
)

type userPostgresRepository struct {
	db *gorm.DB
}

func NewUserPostgresRepository(db *gorm.DB) UserRepository {
	return &userPostgresRepository{db: db}
}

func (r *userPostgresRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userPostgresRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

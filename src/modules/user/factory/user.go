package factory

import (
	"github.com/kev1nandreas/go-rest-api-template/src/modules/user/repository"
	"github.com/kev1nandreas/go-rest-api-template/src/modules/user/usecase"
	"gorm.io/gorm"
)

func NewUserUsecase(db *gorm.DB) usecase.UserUsecase {
	repo := repository.NewUserPostgresRepository(db)
	return usecase.NewUserUsecase(repo)
}

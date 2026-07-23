package factory

import (
	"github.com/kev1nandreas/go-rest-api-template/src/modules/book/repository"
	"github.com/kev1nandreas/go-rest-api-template/src/modules/book/usecase"
	"gorm.io/gorm"
)

func NewBookUsecase(db *gorm.DB) usecase.BookUsecase {
	repo := repository.NewBookPostgresRepository(db)
	return usecase.NewBookUsecase(repo)
}

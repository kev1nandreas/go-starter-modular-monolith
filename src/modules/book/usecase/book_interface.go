package usecase

import (
	"github.com/kev1nandreas/go-rest-api-template/src/modules/book/dto"
	"github.com/kev1nandreas/go-rest-api-template/src/modules/book/entity"
)

type BookUsecase interface {
	FindAll(offset, limit int) ([]entity.Book, error)
	FindByID(id string) (*entity.Book, error)
	Create(req dto.CreateBookRequest) (*entity.Book, error)
	Update(id string, req dto.UpdateBookRequest) (*entity.Book, error)
	Delete(id string) error
}

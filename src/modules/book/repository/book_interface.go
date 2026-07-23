package repository

import (
	"github.com/kev1nandreas/go-rest-api-template/src/modules/book/entity"
)

type BookRepository interface {
	FindAll(offset, limit int) ([]entity.Book, error)
	FindByID(id string) (*entity.Book, error)
	Create(book *entity.Book) error
	Update(book *entity.Book) error
	Delete(book *entity.Book) error
}

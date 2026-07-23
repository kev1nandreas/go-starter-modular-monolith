package repository

import (
	"github.com/kev1nandreas/go-rest-api-template/src/modules/book/entity"
	"gorm.io/gorm"
)

type bookPostgresRepository struct {
	db *gorm.DB
}

func NewBookPostgresRepository(db *gorm.DB) BookRepository {
	return &bookPostgresRepository{db: db}
}

func (r *bookPostgresRepository) FindAll(offset, limit int) ([]entity.Book, error) {
	var books []entity.Book
	if err := r.db.Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookPostgresRepository) FindByID(id string) (*entity.Book, error) {
	var book entity.Book
	if err := r.db.Where("id = ?", id).First(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookPostgresRepository) Create(book *entity.Book) error {
	return r.db.Create(book).Error
}

func (r *bookPostgresRepository) Update(book *entity.Book) error {
	return r.db.Model(book).Updates(entity.Book{
		Title:  book.Title,
		Author: book.Author,
	}).Error
}

func (r *bookPostgresRepository) Delete(book *entity.Book) error {
	return r.db.Delete(book).Error
}

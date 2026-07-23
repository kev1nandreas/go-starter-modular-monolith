package usecase

import (
	"github.com/kev1nandreas/go-rest-api-template/src/modules/book/dto"
	"github.com/kev1nandreas/go-rest-api-template/src/modules/book/entity"
	"github.com/kev1nandreas/go-rest-api-template/src/modules/book/repository"
)

type bookUsecase struct {
	repo repository.BookRepository
}

func NewBookUsecase(repo repository.BookRepository) BookUsecase {
	return &bookUsecase{repo: repo}
}

func (u *bookUsecase) FindAll(offset, limit int) ([]entity.Book, error) {
	return u.repo.FindAll(offset, limit)
}

func (u *bookUsecase) FindByID(id string) (*entity.Book, error) {
	return u.repo.FindByID(id)
}

func (u *bookUsecase) Create(req dto.CreateBookRequest) (*entity.Book, error) {
	book := &entity.Book{
		Title:  req.Title,
		Author: req.Author,
	}

	if err := u.repo.Create(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (u *bookUsecase) Update(id string, req dto.UpdateBookRequest) (*entity.Book, error) {
	book, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}

	if err := u.repo.Update(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (u *bookUsecase) Delete(id string) error {
	book, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}

	return u.repo.Delete(book)
}

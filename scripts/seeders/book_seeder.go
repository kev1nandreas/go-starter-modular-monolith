package seeders

import (
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/kev1nandreas/go-rest-api-template/src/modules/book/entity"
	"gorm.io/gorm"
)

func SeedBooks(db *gorm.DB, count int) error {
	gofakeit.Seed(0)

	for i := 0; i < count; i++ {
		book := entity.Book{
			ID:     uuid.New(),
			Title:  gofakeit.BookTitle(),
			Author: gofakeit.BookAuthor(),
		}

		if err := db.Create(&book).Error; err != nil {
			log.Printf("Failed to create book: %v", err)
			return err
		}
	}

	log.Printf("Successfully seeded %d books", count)
	return nil
}

func ClearBooks(db *gorm.DB) error {
	if err := db.Exec("DELETE FROM books").Error; err != nil {
		return err
	}
	log.Println("Successfully cleared books table")
	return nil
}

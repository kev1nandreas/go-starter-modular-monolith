package factory

import (
	"fmt"
	"log"
	"time"

	"github.com/kev1nandreas/go-rest-api-template/config"
	bookEntity "github.com/kev1nandreas/go-rest-api-template/src/modules/book/entity"
	userEntity "github.com/kev1nandreas/go-rest-api-template/src/modules/user/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	isUseMigration := config.GetEnvBool("APP_MIGRATIONS", false)
	var database *gorm.DB
	var err error

	dbHostname := config.GetEnvString("POSTGRES_HOST", "localhost")
	dbName := config.GetEnvString("POSTGRES_DB", "postgres")
	dbUser := config.GetEnvString("POSTGRES_USER", "postgres")
	dbPass := config.GetEnvString("POSTGRES_PASSWORD", "password")
	dbPort := config.GetEnvString("POSTGRES_PORT", "5432")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHostname,
		dbUser,
		dbPass,
		dbName,
		dbPort,
	)

	for i := 1; i <= 3; i++ {
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		} else {
			log.Printf("Attempt %d: Failed to initialize database. Retrying...", i)
			time.Sleep(3 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after 3 attempts: %v", err)
	}

	if isUseMigration {
		if err := Migrate(database); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
	}

	return database
}

func Migrate(db *gorm.DB) error {
	// Enable UUID extension for PostgreSQL
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	return db.AutoMigrate(&bookEntity.Book{}, &userEntity.User{})
}

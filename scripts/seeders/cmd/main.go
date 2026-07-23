package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kev1nandreas/go-rest-api-template/config"
	"github.com/kev1nandreas/go-rest-api-template/scripts/seeders"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a seeder action: 'up' or 'down'")
	}

	action := os.Args[1]

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	db := connectDatabase()

	switch action {
	case "up":
		log.Println("Running database seeders (up)...")
		if err := seeders.RunAllSeeders(db); err != nil {
			log.Fatalf("Failed to run seeders: %v", err)
		}
		log.Println("Seeders completed successfully!")
	case "down":
		log.Println("Clearing all seeded data (down)...")
		if err := seeders.ClearAllData(db); err != nil {
			log.Fatalf("Failed to clear data: %v", err)
		}
		log.Println("Data cleared successfully!")
	default:
		log.Fatalf("Unknown action: %s. Use 'up' or 'down'", action)
	}
}

func connectDatabase() *gorm.DB {
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
			log.Printf("Attempt %d: Failed to connect to database. Retrying...", i)
			time.Sleep(3 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after 3 attempts: %v", err)
	}

	log.Println("Successfully connected to database")
	return database
}

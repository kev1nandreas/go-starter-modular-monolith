package seeders

import (
	"log"

	"gorm.io/gorm"
)

func RunAllSeeders(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	if err := SeedUsers(db, 10); err != nil {
		return err
	}

	if err := SeedBooks(db, 20); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func ClearAllData(db *gorm.DB) error {
	log.Println("Clearing all seeded data...")

	if err := ClearBooks(db); err != nil {
		return err
	}

	if err := ClearUsers(db); err != nil {
		return err
	}

	log.Println("All data cleared successfully!")
	return nil
}

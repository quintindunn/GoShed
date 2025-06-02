package database

import (
	"com.quintindev/WebShed/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=5678 dbname=WebShed port=5432 sslmode=disable"
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

func AutoMigrations() {
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("AutoMigrate User failed:", err)
	}

	if err := DB.AutoMigrate(&models.RollingCode{}); err != nil {
		log.Fatal("AutoMigrate RollingCode failed:", err)
	}

	if err := DB.AutoMigrate(&models.AllocatedCode{}); err != nil {
		log.Fatal("AutoMigrate AllocatedCode failed:", err)
	}
}

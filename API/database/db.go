package database

import (
	"com.quintindev/APIShed/models"
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
	if err := DB.AutoMigrate(&models.Config{}); err != nil {
		log.Fatal("AutoMigrate Config failed:", err)
	}

	var count int64
	DB.Table("configs").Count(&count)
	if count == 0 {
		log.Println("No config found, creating...")
		DB.Exec("INSERT INTO configs DEFAULT VALUES")
		DB.Table("configs").Count(&count)
	}

	if count > 1 {
		log.Fatal("Too many configs!")
	}

	if count == 0 {
		log.Fatal("Config creation failed!")
	}
}

package utils

import (
	"com.quintindev/WebShed/database"
	"log"
)

func QueryConfigValue[T any](key string) T {
	var value T

	err := database.DB.Table("configs").
		Select(key).
		Limit(1).
		Scan(&value).Error

	if err != nil {
		log.Fatalf("Error getting key %q from table \"configs\": %+v", key, err)
	}

	return value
}

func SetConfigValue[T any](column string, value T) {
	// Fetch the first row to get a primary key or full row
	var row map[string]interface{}
	err := database.DB.Table("configs").
		Select("id"). // or whatever the primary key is
		Limit(1).
		Find(&row).Error

	if err != nil {
		log.Fatalf("Error fetching first config row: %+v", err)
	}
	if row["id"] == nil {
		log.Fatalf("No rows found in configs table")
	}

	// Update the desired column in that row
	err = database.DB.Table("configs").
		Where("id = ?", row["id"]).
		Update(column, value).Error

	if err != nil {
		log.Fatalf("Error updating column %q: %+v", column, err)
	}
}

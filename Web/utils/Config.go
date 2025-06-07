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

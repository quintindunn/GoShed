package main

import (
	"com.quintindev/WebShed/config"
	"com.quintindev/WebShed/database"
	"com.quintindev/WebShed/models"
	"com.quintindev/WebShed/routes"
	"log"
)

func main() {
	cfg := config.Load()

	database.Init()
	if err := database.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}
	router := routes.SetupRouter()

	err := router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal("Router run failed:", err)

	}
}

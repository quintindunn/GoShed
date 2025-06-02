package main

import (
	"com.quintindev/WebShed/config"
	"com.quintindev/WebShed/database"
	"com.quintindev/WebShed/models"
	"com.quintindev/WebShed/routes"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Starting WebShed.")

	fmt.Println("Loading config.")
	cfg := config.Load()

	fmt.Println("Initializing and migrating database.")
	database.Init()
	if err := database.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	fmt.Println("Setting up router.")
	router := routes.SetupRouter()

	fmt.Println("Starting Server on port: " + cfg.Port)
	err := router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal("Router run failed:", err)

	}
}

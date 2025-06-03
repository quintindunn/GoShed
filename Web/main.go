package main

import (
	"com.quintindev/WebShed/config"
	"com.quintindev/WebShed/database"
	"com.quintindev/WebShed/routes"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Starting APIShed.")

	fmt.Println("Loading config.")
	cfg := config.Load()

	fmt.Println("Initializing and migrating database.")
	database.Init()
	database.AutoMigrations()

	fmt.Println("Setting up router.")
	router := routes.SetupRouter()

	fmt.Println("Starting Server on port: " + cfg.Port)
	err := router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal("Router run failed:", err)

	}
}

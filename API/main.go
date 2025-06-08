package main

import (
	"com.quintindev/APIShed/config"
	"com.quintindev/APIShed/database"
	"com.quintindev/APIShed/hardware"
	"com.quintindev/APIShed/routes"
	"com.quintindev/APIShed/util"
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

	hardware.SetLockState(util.QueryConfigValue[bool]("lock_state"))

	fmt.Println("Starting Server on port: " + cfg.Port)
	err := router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal("Router run failed:", err)

	}
}

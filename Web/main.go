package main

import (
	"com.quintindev/WebShed/config"
	"com.quintindev/WebShed/routes"
)

func main() {
	cfg := config.Load()

	router := routes.SetupRouter()

	router.Run(":" + cfg.Port)
}

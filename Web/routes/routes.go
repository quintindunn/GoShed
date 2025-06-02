package routes

import (
	"com.quintindev/WebShed/controllers"
	"com.quintindev/WebShed/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", controllers.Page)
	r.GET("/api", controllers.Home)

	return r
}

package routes

import (
	"com.quintindev/APIShed/controllers"
	"com.quintindev/APIShed/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	r.POST("/api/setlock", controllers.SetLock)
	r.GET("/api/getlocked", controllers.GetLocked)

	return r
}

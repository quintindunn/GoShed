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

	r.GET("/login", controllers.LoginPage)
	r.POST("/login", controllers.LoginSubmit)

	r.GET("/register", controllers.RegisterPage)
	r.POST("/register", controllers.Register)

	r.POST("/api/login", controllers.LoginAPI)
	r.POST("/api/register", controllers.Register)

	r.GET("/logout", controllers.Logout)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired())
	{
		authorized.GET("/", controllers.Page)
		authorized.GET("/users", controllers.ListUsers)
		authorized.GET("/locks", controllers.Lock)
		authorized.POST("/api/lock", controllers.SetLockAPI)
		authorized.POST("/api/refreshCards", controllers.ResetRollingCodesAPI)
		authorized.POST("/api/addUserCode", controllers.AddUserCodeAPI)
	}

	return r
}

package controllers

import (
	"com.quintindev/WebShed/utils"
	"github.com/gin-gonic/gin"
)

func Lock(c *gin.Context) {
	data := gin.H{}
	utils.Render(c, 200, "locks", data)
}

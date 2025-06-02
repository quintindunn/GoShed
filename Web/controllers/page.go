package controllers

import "github.com/gin-gonic/gin"

func Page(c *gin.Context) {
	data := gin.H{}
	c.HTML(200, "home", data)
}

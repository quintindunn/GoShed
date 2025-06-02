package controllers

import (
	"net/http"

	"com.quintindev/WebShed/database"
	"com.quintindev/WebShed/models"
	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

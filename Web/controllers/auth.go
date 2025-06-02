package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"com.quintindev/WebShed/config"
	"com.quintindev/WebShed/database"
	"com.quintindev/WebShed/models"
)

type RegisterPayload struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func RegisterPage(c *gin.Context) {
	if token, err := c.Cookie("jwt_token"); err == nil && token != "" {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	c.HTML(http.StatusOK, "register", gin.H{"Error": ""})
}

func Register(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	if name == "" || email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "register", gin.H{"Error": "All fields are required"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register", gin.H{"Error": "Server error"})
		return
	}

	user := models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.HTML(http.StatusBadRequest, "register", gin.H{"Error": "Email already in use"})
			return
		}
		c.HTML(http.StatusInternalServerError, "register", gin.H{"Error": "Could not create account"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func LoginPage(c *gin.Context) {
	if token, err := c.Cookie("jwt_token"); err == nil && token != "" {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	c.HTML(http.StatusOK, "login", gin.H{"Error": ""})
}

func LoginSubmit(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.HTML(http.StatusUnauthorized, "login", gin.H{"Error": "Invalid email or password"})
			return
		}
		c.HTML(http.StatusInternalServerError, "login", gin.H{"Error": "Server error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login", gin.H{"Error": "Invalid email or password"})
		return
	}

	tokenString, err := config.GenerateToken(user.ID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login", gin.H{"Error": "Could not generate token"})
		return
	}

	c.SetCookie("jwt_token", tokenString, 3600*24, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func LoginAPI(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	tokenString, err := config.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Logout(c *gin.Context) {
	c.SetCookie("jwt_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/login")
}

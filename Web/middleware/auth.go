package middleware

import (
	"com.quintindev/WebShed/database"
	"com.quintindev/WebShed/models"
	"net/http"
	"strings"

	"com.quintindev/WebShed/config"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tokenCookie, err := c.Cookie("jwt_token"); err == nil && tokenCookie != "" {
			if claims, err := config.ParseToken(tokenCookie); err == nil {
				userID := uint(claims["user_id"].(float64))

				var user models.User
				if err := database.DB.First(&user, userID).Error; err == nil {
					c.Set("currentUser", user)
					c.Next()
					return
				}
			}
		}

		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			if claims, err := config.ParseToken(parts[1]); err == nil {
				userID := uint(claims["user_id"].(float64))
				var user models.User
				if err := database.DB.First(&user, userID).Error; err == nil {
					c.Set("currentUser", user)
					c.Next()
					return
				}
			}
		}

		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
	}
}

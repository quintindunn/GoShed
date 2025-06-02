package middleware

import (
	"net/http"
	"strings"

	"com.quintindev/WebShed/config"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tokenCookie, err := c.Cookie("jwt_token"); err == nil && tokenCookie != "" {
			if claims, err := config.ParseToken(tokenCookie); err == nil {
				c.Set("user_id", uint(claims["user_id"].(float64)))
				c.Next()
				return
			}
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := config.ParseToken(tokenString)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Next()
	}
}

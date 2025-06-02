package utils

import (
	"com.quintindev/WebShed/models"
	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, code int, tmpl string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	if u, exists := c.Get("currentUser"); exists {
		if user, ok := u.(models.User); ok {
			data["CurrentUser"] = user
		} else {
		}
	} else {
	}
	c.HTML(code, tmpl, data)
}

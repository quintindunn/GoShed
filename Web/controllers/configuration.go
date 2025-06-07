package controllers

import (
	"com.quintindev/WebShed/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Configuration(c *gin.Context) {
	utils.Render(c, http.StatusOK, "configuration", gin.H{
		"needAdminPin": utils.QueryConfigValue[bool]("need_admin_pin_for_user_management"),
		"unlockTime":   utils.QueryConfigValue[int64]("unlock_time"),
	})
}

package controllers

import (
	"com.quintindev/WebShed/database"
	"com.quintindev/WebShed/hardwareInterface"
	"com.quintindev/WebShed/models"
	"com.quintindev/WebShed/utils"
	"github.com/gin-gonic/gin"
	"time"
)

func Page(c *gin.Context) {
	hardwareInterface.GetExpiredCodes()

	var ModelRollingCodes []models.RollingCode
	database.DB.Find(&ModelRollingCodes, "nullified = ?", false)

	var RollingCodes []gin.H
	for _, code := range ModelRollingCodes {
		RollingCodes = append(RollingCodes, gin.H{
			"code":   code.Code,
			"expiry": time.Unix(code.Expiry, 0).Format("01-02-06 3:04 PM"),
		})
	}
	data := gin.H{
		"isLocked":     hardwareInterface.GetGetLocked(),
		"rollingCodes": RollingCodes,
	}
	utils.Render(c, 200, "home", data)
}

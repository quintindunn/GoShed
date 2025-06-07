package controllers

import (
	"com.quintindev/APIShed/util"
	_ "com.quintindev/APIShed/util"
	"github.com/gin-gonic/gin"
)

func ExpireOldCodes(c *gin.Context) {
	rollingCodesChanged := util.UpdateExpiredRollingCodes()
	allocatedCodesChanged := util.NullifyAllocatedCodes()

	c.JSON(200, gin.H{
		"rollingCodesChanged":   rollingCodesChanged,
		"allocatedCodesChanged": allocatedCodesChanged,
	})
}

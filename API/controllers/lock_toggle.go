package controllers

import (
	"com.quintindev/APIShed/hardware"

	"github.com/gin-gonic/gin"
	"net/http"
)

type SetLockRequest struct {
	State bool `json:"state"`
}

func SetLock(c *gin.Context) {
	var json SetLockRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
	}

	hardware.SetLockState(json.State)

	c.JSON(http.StatusOK, gin.H{
		"newState": hardware.LockState.Locked,
	})
}

func GetLocked(c *gin.Context) {
	isLocked := hardware.GetLockedState()

	c.JSON(http.StatusOK, gin.H{
		"state": isLocked,
	})
}

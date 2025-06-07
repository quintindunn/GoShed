package controllers

import (
	"com.quintindev/APIShed/audit"
	"com.quintindev/APIShed/database"
	"com.quintindev/APIShed/hardware"
	"com.quintindev/APIShed/models"

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

type UnlockRequest struct {
	Code      string `json:"code"`
	Initiator string `json:"initiator"`
}

func AttemptUnlock(c *gin.Context) {
	var json UnlockRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
		})
		return
	}

	var allocatedCodes []models.AllocatedCode
	database.DB.Model(&models.AllocatedCode{}).
		Where("nullified = ?", false).
		Where("code = ?", json.Code).Find(&allocatedCodes)

	var rollingCodes []models.RollingCode
	database.DB.Model(&models.RollingCode{}).
		Where("nullified = ?", false).
		Where("code = ?", json.Code).Find(&rollingCodes)

	failure := len(allocatedCodes) == 0 && len(rollingCodes) == 0
	if failure {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
	})

	if len(allocatedCodes) != 0 {
		audit.UnlockByAllocatedCode(allocatedCodes[0], json.Initiator)
		hardware.HandleCodedUnlock()
		return
	} else if len(rollingCodes) != 0 {
		audit.UnlockByRollingCode(rollingCodes[0], json.Initiator)
		hardware.HandleCodedUnlock()
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"success": false,
		})
	}
}

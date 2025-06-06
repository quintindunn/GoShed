package controllers

import (
	"com.quintindev/APIShed/audit"
	"com.quintindev/APIShed/database"
	"com.quintindev/APIShed/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"time"
)

func ExpireOldCodes(c *gin.Context) {
	rollingCodesChanged := updateExpiredRollingCodes()
	allocatedCodesChanged := nullifyAllocatedCodes()

	c.JSON(200, gin.H{
		"rollingCodesChanged":   rollingCodesChanged,
		"allocatedCodesChanged": allocatedCodesChanged,
	})
}

func updateExpiredRollingCodes() int {
	unix := time.Now().Unix()

	var expiredModels []models.RollingCode
	database.DB.Model(&models.RollingCode{}).
		Where("nullified = ?", false).
		Where("expiry < ?", unix).Find(&expiredModels)

	nullifiedRollingCodesCount := len(expiredModels)
	var nullifiedCodes []string
	for _, model := range expiredModels {
		database.DB.Model(&model).Update("nullified", true)
		nullifiedCodes = append(nullifiedCodes, model.Code)
	}

	now := time.Now()
	var ModelRollingCodes []models.RollingCode
	var newCodes []string

	for i := 0; i < nullifiedRollingCodesCount; i++ {
		num := rand.Intn(1000000)       // 0 to 999999
		str := fmt.Sprintf("%06d", num) // zero-padded string
		ModelRollingCodes = append(ModelRollingCodes, models.RollingCode{
			Code:      str,
			Expiry:    now.Unix() + 86400,
			Nullified: false,
		})
		newCodes = append(newCodes, str)
	}

	if len(ModelRollingCodes) == 0 {
		return 0
	}

	audit.NullifyRollingCodes(nullifiedCodes)
	audit.CreateNewRollingCodes(newCodes)

	if err := database.DB.Create(&ModelRollingCodes).Error; err != nil {
		log.Println("Failed to insert rolling codes:", err)
		return -1
	} else {
		return nullifiedRollingCodesCount
	}
}

func nullifyAllocatedCodes() int {
	unix := time.Now().Unix()

	var expiredModels []models.AllocatedCode
	database.DB.Model(&models.AllocatedCode{}).
		Where("nullified = ?", false).
		Where("expiry < ?", unix).Find(&expiredModels)

	nullifiedAllocatedCodesCount := len(expiredModels)

	for _, model := range expiredModels {
		database.DB.Model(&model).Update("nullified", true)
	}

	return nullifiedAllocatedCodesCount
}

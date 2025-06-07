package util

import (
	"com.quintindev/APIShed/audit"
	"com.quintindev/APIShed/database"
	"com.quintindev/APIShed/models"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func UpdateExpiredRollingCodes() int {
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

func NullifyAllocatedCodes() int {
	unix := time.Now().Unix()

	var expiredModels []models.AllocatedCode
	database.DB.Model(&models.AllocatedCode{}).
		Where("nullified = ?", false).
		Where("expiry < ?", unix).Find(&expiredModels)

	nullifiedAllocatedCodesCount := len(expiredModels)

	for _, model := range expiredModels {
		database.DB.Model(&model).Update("nullified", true)
	}

	var nullifiedCodes [][]string

	for _, model := range expiredModels {
		nullifiedCodes = append(nullifiedCodes, []string{model.Name, model.Code})
		database.DB.Model(&model).Update("nullified", true)
	}

	if len(nullifiedCodes) != 0 {
		audit.NullifyAllocatedCodes(nullifiedCodes)
	}

	return nullifiedAllocatedCodesCount
}

func QueryConfigInt64(key string) int64 {
	var value int64

	err := database.DB.Table("configs").
		Select(key).
		Limit(1).
		Scan(&value).Error

	if err != nil {
		log.Fatalf("Error getting key \"%s\" from table \"configs\": %+v", key, err)
	}

	return value
}

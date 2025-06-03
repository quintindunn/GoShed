package main

import (
	"com.quintindev/CronShed/database"
	"com.quintindev/CronShed/models"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	database.Init()
	database.AutoMigrations()
	updateExpiredRollingCodes()
	nullifyAllocatedCodes()
}

func updateExpiredRollingCodes() {
	unix := time.Now().Unix()

	var expiredModels []models.RollingCode
	database.DB.Model(&models.RollingCode{}).
		Where("nullified = ?", false).
		Where("expiry < ?", unix).Find(&expiredModels)

	nullifiedRollingCodesCount := len(expiredModels)

	for _, model := range expiredModels {
		database.DB.Model(&model).Update("nullified", true)
	}

	now := time.Now()
	var ModelRollingCodes []models.RollingCode

	for i := 0; i < nullifiedRollingCodesCount; i++ {
		num := rand.Intn(1000000)       // 0 to 999999
		str := fmt.Sprintf("%06d", num) // zero-padded string
		ModelRollingCodes = append(ModelRollingCodes, models.RollingCode{
			Code:      str,
			Expiry:    now.Unix() + 86400,
			Nullified: false,
		})
	}

	if len(ModelRollingCodes) == 0 {
		log.Println("No rolling codes expired!")
		return
	}

	if err := database.DB.Create(&ModelRollingCodes).Error; err != nil {
		log.Println("Failed to insert rolling codes:", err)
	} else {
		log.Printf("Inserted %d rolling codes!\n", nullifiedRollingCodesCount)
	}
}

func nullifyAllocatedCodes() {
	unix := time.Now().Unix()

	var expiredModels []models.AllocatedCode
	database.DB.Model(&models.AllocatedCode{}).
		Where("nullified = ?", false).
		Where("expiry < ?", unix).Find(&expiredModels)

	nullifiedAllocatedCodesCount := len(expiredModels)

	for _, model := range expiredModels {
		database.DB.Model(&model).Update("nullified", true)
	}

	if nullifiedAllocatedCodesCount == 0 {
		log.Println("No allocated codes expired!")
	} else {
		log.Printf("Nullified %d allocated codes!\n", nullifiedAllocatedCodesCount)
	}
}

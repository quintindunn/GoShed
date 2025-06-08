package audit

import (
	"com.quintindev/WebShed/database"
	"com.quintindev/WebShed/models"
	"fmt"
	"time"
)

func Log(msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s | %s\n", timestamp, msg)

	log := models.Log{
		Msg: msg,
	}

	if err := database.DB.Create(&log).Error; err != nil {
		fmt.Println("Error creating log:", err)
	}
}

func LogRefreshRollingCodes(newCodes []string) {
	msg := fmt.Sprintf("SYSTEM | Refreshing %d new codes: %+v", len(newCodes), newCodes)
	Log(msg)
}

func LogAddNewCode(name string, newCode string, expiration int64) {
	msg := fmt.Sprintf("SYSTEM | Adding new code \"%s\" for user \"%s\" that expires %d!", newCode, name, expiration)
	Log(msg)
}

func LogRemoveCode(targetUUID string) {
	var allocatedCodes []models.AllocatedCode
	database.DB.Find(&allocatedCodes, "uuid = ?", targetUUID)

	if len(allocatedCodes) < 0 {
		Log("SYSTEM ERROR | Couldn't find allocated code that was removed!")
	}

	msg := fmt.Sprintf("SYSTEM | Removing code \"%s\"!", allocatedCodes[0].Code)
	Log(msg)
}

func LogNewConfiguration(configJsonString string) {
	msg := fmt.Sprintf("SYSTEM | New Config: %+v", configJsonString)
	Log(msg)
}
